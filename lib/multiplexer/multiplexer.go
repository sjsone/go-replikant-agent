package multiplexer

import (
	"context"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
)

type Multiplexer interface {
	GetActiveDirectivesForContext(ctx context.Context, ac agentic_context.AgentContext) []directive.Directive
	GetAllDirectives() []directive.Directive
}
