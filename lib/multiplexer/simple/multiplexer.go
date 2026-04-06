package multiplexer

import (
	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
)

type SimpleMultiplexer struct {
	directives []directive.Directive
}

func NewSimpleMultiplexer(directives []directive.Directive) *SimpleMultiplexer {
	return &SimpleMultiplexer{
		directives: directives,
	}
}

func (m *SimpleMultiplexer) GetActiveDirectivesForContext(ac agentic_context.AgentContext) []directive.Directive {
	return m.directives
}

// GetAllDirectives returns all directives available to this multiplexer.
func (m *SimpleMultiplexer) GetAllDirectives() []directive.Directive {
	return m.directives
}
