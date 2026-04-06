package session

import (
	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// SessionDelegate observes and reacts to events during AgenticSession execution.
// All methods are optional - implementations can choose which to override.
type SessionDelegate interface {
	// SessionOnPartAdded is called when a new context part is added to the conversation.
	SessionOnPartAdded(part *agentic_context.ContextPart)

	// SessionOnToolCallsReceived is called when the LLM returns tool calls.
	SessionOnToolCallsReceived(calls []tool.FunctionCall)

	// SessionOnToolExecuted is called after a tool is executed.
	SessionOnToolExecuted(call tool.FunctionCall, result tool.FunctionResult)

	// SessionOnStreamingChunk is called for each chunk of streaming content from the LLM.
	SessionOnStreamingChunk(chunk string)

	// SessionOnRequestSent is called before sending a request to the LLM.
	SessionOnRequestSent(messages []connector.Message, directives []directive.Directive)

	// SessionOnLoopIteration is called at the start of each loop iteration.
	SessionOnLoopIteration(iteration int)

	// SessionOnLoopEnd is called when the agent loop finishes (or is broken out of).
	SessionOnLoopEnd()
}
