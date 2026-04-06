package loop

import "github.com/sjsone/go-replikant-agent/lib/agentic_context"

type LoopStopReason = string

// Reason constants describe why the loop controller made its decision.
const (
	ReasonNilContext          LoopStopReason = "nil_context"
	ReasonRepetitiveToolCalls LoopStopReason = "repetitive_tool_calls"
	ReasonCancelled           LoopStopReason = "cancelled"
	ReasonStopFlag            LoopStopReason = "stop_flag"
	ReasonUserSource          LoopStopReason = "user_source"
	ReasonToolUse             LoopStopReason = "tool_use"
	ReasonToolResult          LoopStopReason = "tool_result"
	ReasonSystemSource        LoopStopReason = "system_source"
	ReasonDefaultStop         LoopStopReason = "default_stop"
)

// LoopDecision describes a single loop continuation decision.
type LoopDecision struct {
	Continue bool
	Reason   LoopStopReason
	Part     *agentic_context.ContextPart
}

// Delegate observes loop control decisions without influencing them.
type Delegate interface {
	LoopOnLoopDecision(decision LoopDecision)
}
