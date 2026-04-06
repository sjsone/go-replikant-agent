package loop

import "github.com/sjsone/go-replikant-agent/lib/agentic_context"

// LoopController determines whether the agent loop should continue.
type LoopController interface {
	LoopAgain(agentContext *agentic_context.AgentContext) bool
}
