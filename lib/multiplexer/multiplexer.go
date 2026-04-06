package multiplexer

import (
	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
)

type Multiplexer interface {
	GetActiveDirectivesForContext(agentic_context.AgentContext) []directive.Directive
	GetAllDirectives() []directive.Directive
}
