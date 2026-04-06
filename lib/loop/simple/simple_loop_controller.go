package simple

import (
	"fmt"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/loop"
)

const (
	defaultRepetitiveToolCallWindow    = 10
	defaultRepetitiveToolCallThreshold = 5
)

// SimpleLoopController is the default implementation of LoopController.
type SimpleLoopController struct {
	// Delegate receives notifications about loop decisions. Nil-safe.
	delegate loop.Delegate

	// RepetitiveToolCallWindow is the number of recent tool calls to inspect.
	// Zero means 10.
	RepetitiveToolCallWindow int

	// RepetitiveToolCallThreshold is the number of identical tool calls (same name + args)
	// within the window that triggers loop termination. Zero means 5.
	RepetitiveToolCallThreshold int
}

func NewSimpleLoopController() *SimpleLoopController {
	return &SimpleLoopController{
		RepetitiveToolCallWindow:    defaultRepetitiveToolCallWindow,
		RepetitiveToolCallThreshold: defaultRepetitiveToolCallThreshold,
	}
}

func (c *SimpleLoopController) SetDelegate(delegate *loop.Delegate) {
	c.delegate = *delegate
}

func (c *SimpleLoopController) LoopAgain(agentContext *agentic_context.AgentContext) bool {
	decision := c.evaluate(agentContext)
	c.notifyDelegate(decision)
	return decision.Continue
}

// evaluate returns the loop decision for the current context without side effects.
func (c *SimpleLoopController) evaluate(agentContext *agentic_context.AgentContext) loop.LoopDecision {
	part := agentContext.GetLatestPart()
	if part == nil {
		return loop.LoopDecision{Continue: true, Reason: loop.ReasonNilContext, Part: nil}
	}

	// Break out of repetitive tool call loops
	if c.isRepetitiveToolCalls(agentContext.Parts) {
		return loop.LoopDecision{Continue: false, Reason: loop.ReasonRepetitiveToolCalls, Part: part}
	}

	// Cancelled parts stop the loop (operation was interrupted)
	if part.Cancelled {
		return loop.LoopDecision{Continue: false, Reason: loop.ReasonCancelled, Part: part}
	}

	if part.Stop {
		return loop.LoopDecision{Continue: false, Reason: loop.ReasonStopFlag, Part: part}
	}

	if part.Source.IsUser() {
		return loop.LoopDecision{Continue: true, Reason: loop.ReasonUserSource, Part: part}
	}

	if part.ToolUse {
		return loop.LoopDecision{Continue: true, Reason: loop.ReasonToolUse, Part: part}
	}

	if part.IsToolResult() {
		return loop.LoopDecision{Continue: true, Reason: loop.ReasonToolResult, Part: part}
	}

	if part.Source.IsSystem() {
		return loop.LoopDecision{Continue: false, Reason: loop.ReasonSystemSource, Part: part}
	}

	return loop.LoopDecision{Continue: false, Reason: loop.ReasonDefaultStop, Part: part}
}

// notifyDelegate calls the delegate if one is set.
func (c *SimpleLoopController) notifyDelegate(decision loop.LoopDecision) {
	if c.delegate != nil {
		c.delegate.LoopOnLoopDecision(decision)
	}
}

// isRepetitiveToolCalls checks whether any unique tool call (same Name + Arguments)
// appears at least RepetitiveToolCallThreshold times within the last RepetitiveToolCallWindow calls.
func (c *SimpleLoopController) isRepetitiveToolCalls(parts []*agentic_context.ContextPart) bool {
	window := c.RepetitiveToolCallWindow
	threshold := c.RepetitiveToolCallThreshold

	counts := make(map[string]int)
	collected := 0

	for i := len(parts) - 1; i >= 0 && collected < window; i-- {
		p := parts[i]
		if !p.ToolUse {
			continue
		}
		for _, tc := range p.ToolCalls {
			key := fmt.Sprintf("%s:%v", tc.Name, tc.Arguments)
			counts[key]++
			if counts[key] >= threshold {
				return true
			}
			collected++
			if collected >= window {
				return false
			}
		}
	}
	return false
}
