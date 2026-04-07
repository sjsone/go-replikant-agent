package session

import (
	"context"
	"fmt"
	"sync"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/loop"
	"github.com/sjsone/go-replikant-agent/lib/multiplexer"
	prompt_builder "github.com/sjsone/go-replikant-agent/lib/prompt/builder"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type AgenticSession struct {
	directiveMultiplexer multiplexer.Multiplexer
	currentContext       agentic_context.AgentContext
	promptBuilder        prompt_builder.PromptBuilder
	loopController       loop.LoopController
	connector            connector.Connector
	delegate             SessionDelegate
	ctx                  context.Context
	ctxMu                sync.RWMutex
	maxAgentLoopCount    int
}

func NewAgenticSession(directiveMultiplexer multiplexer.Multiplexer, currentContext agentic_context.AgentContext, promptBuilder prompt_builder.PromptBuilder, loopController loop.LoopController, connector connector.Connector) *AgenticSession {
	return &AgenticSession{
		directiveMultiplexer: directiveMultiplexer,
		currentContext:       currentContext,
		promptBuilder:        promptBuilder,
		loopController:       loopController,
		connector:            connector,
		delegate:             nil,
		maxAgentLoopCount:    100,
	}
}

func (as *AgenticSession) SetDelegate(delegate SessionDelegate) {
	as.delegate = delegate
}

// CurrentContext returns the session's current agent context (with all accumulated parts).
func (as *AgenticSession) CurrentContext() *agentic_context.AgentContext {
	return &as.currentContext
}

func (as *AgenticSession) ProcessContextPart(ctx context.Context, part *agentic_context.ContextPart) error {
	// Store the context and cancel function for external cancellation
	as.ctxMu.Lock()
	as.ctx = ctx
	as.ctxMu.Unlock()
	as.currentContext.AddPart(part)
	if as.delegate != nil {
		as.delegate.SessionOnPartAdded(part)
	}

	directives := as.directiveMultiplexer.GetActiveDirectivesForContext(as.currentContext)

	for i := 0; i < as.maxAgentLoopCount && as.loopController.LoopAgain(&as.currentContext); i++ {
		// Check for cancellation before each iteration
		as.ctxMu.RLock()
		sessionCtx := as.ctx
		as.ctxMu.RUnlock()

		select {
		case <-sessionCtx.Done():
			return sessionCtx.Err()
		default:
			// Continue with the loop
		}
		if as.delegate != nil {
			as.delegate.SessionOnLoopIteration(i)
		}
		var err error
		err = as.loopInner(sessionCtx, directives)
		if err != nil {
			return fmt.Errorf("loop iteration %d: %w", i, err)
		}
	}

	if as.delegate != nil {
		as.delegate.SessionOnLoopEnd()
	}

	return nil
}

func (as *AgenticSession) buildSystemMessageForConnector(directives []directive.Directive) connector.Message {
	builtPrompt := as.promptBuilder.Build(directives)
	return connector.NewSystemMessage(builtPrompt.Raw)
}

func (as *AgenticSession) buildMessagesForConnector(directives []directive.Directive) *[]connector.Message {
	// TODO: introduce new method to influence message building (to enable classification-based-context-manipulation)
	// Build messages fresh from current context (system message + all parts)

	systemMessage := as.buildSystemMessageForConnector(directives)
	messages := []connector.Message{systemMessage}

	for _, p := range as.currentContext.Parts {
		messages = append(messages, connector.Message{
			Source: p.Source,
			Text:   p.Raw,
		})
	}

	return &messages
}

func (as *AgenticSession) handleToolCall(ctx context.Context, new_part *agentic_context.ContextPart, directives []directive.Directive) {
	if as.delegate != nil {
		as.delegate.SessionOnToolCallsReceived(new_part.ToolCalls)
	}

	// TODO: move logic into SimpleToolManager / ToolManager interface
	results := make([]tool.FunctionResult, 0)

	for _, call := range new_part.ToolCalls {
		// Find tool callable in directives
		var toolCallable tool.ToolCallable

		for _, d := range directives {
			// Find the callable by name
			for _, tc := range d.GetToolCallables() {
				if tc.GetName() == call.Name {
					toolCallable = tc
					break
				}
			}
			if toolCallable != nil {
				break
			}
		}

		var result tool.FunctionResult
		if toolCallable != nil {
			content, err := toolCallable.Execute(ctx, call.Arguments)
			result = tool.FunctionResult{
				ID:      call.ID,
				Content: content,
				Error:   err != nil,
			}
			if err != nil {
				result.Content = fmt.Sprintf("Error: %v", err)
			}
		} else {
			result = tool.FunctionResult{
				ID:      call.ID,
				Content: fmt.Sprintf("Tool not found: %s", call.Name),
				Error:   true,
			}
		}

		if as.delegate != nil {
			as.delegate.SessionOnToolExecuted(call, result)
		}

		results = append(results, result)
	}

	// Add tool results as new context part
	toolPart := agentic_context.NewToolResultContextPart(results)
	toolPart.ConnectedToolCallContextPart = new_part
	as.currentContext.AddPart(toolPart)
	if as.delegate != nil {
		as.delegate.SessionOnPartAdded(toolPart)
	}
}

// TODO: rename `loopInner` to something better
func (as *AgenticSession) loopInner(ctx context.Context, directives []directive.Directive) error {

	messages := as.buildMessagesForConnector(directives)

	if as.delegate != nil {
		as.delegate.SessionOnRequestSent(*messages, directives)
	}

	err, new_part := as.connector.Send(ctx, messages, directives, func(chunk string) {
		if as.delegate != nil {
			as.delegate.SessionOnStreamingChunk(chunk)
		}
	})

	if err != nil {
		// If the part is non-nil and marked as cancelled, add it to context before returning error
		if new_part != nil && new_part.Cancelled {
			as.currentContext.AddPart(new_part)
			if as.delegate != nil {
				as.delegate.SessionOnPartAdded(new_part)
			}
		}
		return err
	}

	as.currentContext.AddPart(new_part)
	if as.delegate != nil {
		as.delegate.SessionOnPartAdded(new_part)
	}

	// Handle tool execution
	if len(new_part.ToolCalls) > 0 {
		new_part.ConnectedToolCallSourceContextPart = as.currentContext.GetLatestNonToolPart()
		as.handleToolCall(ctx, new_part, directives)
	}

	return nil
}

// Cancel cancels the current session operation by cancelling the stored context.
func (as *AgenticSession) Cancel() {
	as.ctxMu.Lock()
	defer as.ctxMu.Unlock()
}
