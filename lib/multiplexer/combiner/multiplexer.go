package combined

import (
	"context"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	multiplexer_lib "github.com/sjsone/go-replikant-agent/lib/multiplexer"
)

type CombinerMultiplexer struct {
	multiplexers []multiplexer_lib.Multiplexer
}

func NewCombinerMultiplexer(multiplexers []multiplexer_lib.Multiplexer) *CombinerMultiplexer {
	return &CombinerMultiplexer{
		multiplexers: multiplexers,
	}
}

func (m *CombinerMultiplexer) GetActiveDirectivesForContext(ctx context.Context, ac agentic_context.AgentContext) []directive.Directive {
	return m.collectActive(ctx, ac)
}

// GetAllDirectives returns all directives across all multiplexers, deduplicated by name.
func (m *CombinerMultiplexer) GetAllDirectives() []directive.Directive {
	return m.collectAll()
}

// collectActive unions the active directives from each child multiplexer, deduplicating by name.
func (m *CombinerMultiplexer) collectActive(ctx context.Context, ac agentic_context.AgentContext) []directive.Directive {
	seen := make(map[string]struct{})
	var result []directive.Directive

	for _, mx := range m.multiplexers {
		for _, d := range mx.GetActiveDirectivesForContext(ctx, ac) {
			name := d.GetName()
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			result = append(result, d)
		}
	}

	return result
}

// collectAll unions all directives from each child multiplexer, deduplicating by name.
func (m *CombinerMultiplexer) collectAll() []directive.Directive {
	seen := make(map[string]struct{})
	var result []directive.Directive

	for _, mx := range m.multiplexers {
		for _, d := range mx.GetAllDirectives() {
			name := d.GetName()
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			result = append(result, d)
		}
	}

	return result
}
