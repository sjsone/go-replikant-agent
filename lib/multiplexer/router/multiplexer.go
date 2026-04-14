package multiplexer

import (
	"context"
	"log"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/router"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type RouterMultiplexer struct {
	directives          []directive.Directive
	directiveNameMap    map[string]directive.Directive
	router              router.Router
	lastRoutingDecision *router.RoutingDecision
}

// NewRouterMultiplexer creates a new RouterMultiplexer that uses LLM-based routing.
func NewRouterMultiplexer(
	directives []directive.Directive,
	r router.Router,
) *RouterMultiplexer {
	nameMap := make(map[string]directive.Directive, len(directives))
	for _, d := range directives {
		nameMap[d.GetName()] = d
	}

	return &RouterMultiplexer{
		directives:       directives,
		directiveNameMap: nameMap,
		router:           r,
	}
}

// GetActiveDirectivesForContext uses LLM routing to determine active directives.
func (m *RouterMultiplexer) GetActiveDirectivesForContext(ctx context.Context, ac agentic_context.AgentContext) []directive.Directive {
	// Extract latest user message
	userMsg := m.extractUserMessage(ac)
	if userMsg == "" {
		// TODO: handle non-found user message
		return m.directives
	}

	options := m.routingOptionsFromDirectives()

	routingResult, err := m.router.Route(ctx, userMsg, options)
	if err != nil {
		log.Printf("routing failed: %v, falling back to all directives", err)
		return m.directives
	}
	if routingResult == nil {
		return m.directives
	}

	// Store the routing decision
	if routingResult.Decision != nil {
		m.lastRoutingDecision = routingResult.Decision
	}

	result := m.directivesFromRoutingOptions(routingResult.SelectedOptions)

	return result
}

// extractUserMessage extracts the latest user message from the context.
func (m *RouterMultiplexer) extractUserMessage(ctx agentic_context.AgentContext) string {
	for i := len(ctx.Parts) - 1; i >= 0; i-- {
		if ctx.Parts[i].Source.IsUser() {
			return ctx.Parts[i].Raw
		}
	}
	// TODO: do not return an empty string but an error instead
	return ""
}

// routingOptionsFromDirectives converts directives to routing options for the router.
func (m *RouterMultiplexer) routingOptionsFromDirectives() []*router.RoutingOption {
	options := make([]*router.RoutingOption, 0, len(m.directives))

	for _, d := range m.directives {
		text := d.GetPrompt().Raw
		tools := tool.ToolsFromCallables(d.GetToolCallables())

		// Extract tool names
		toolNames := make([]string, 0, len(tools))
		for _, t := range tools {
			toolNames = append(toolNames, t.Name)
		}

		if len(tools) > 0 {
			text += "\n"
			text += "Available tools (do NOT use these as option IDs):\n"
			for _, t := range tools {
				text += "  - " + t.Name + ": " + t.Description
			}
		}

		option := router.RoutingOption{
			Name:      d.GetName(),
			Text:      text,
			ToolNames: toolNames,
		}

		options = append(options, &option)
	}

	return options
}

func (m *RouterMultiplexer) directivesFromRoutingOptions(options []*router.RoutingOption) []directive.Directive {
	directives := make([]directive.Directive, 0, len(options))

	for _, o := range options {
		if d, ok := m.directiveNameMap[o.Name]; ok {
			directives = append(directives, d)
		}
	}

	return directives
}

// GetAllDirectives returns all directives available to this multiplexer.
func (m *RouterMultiplexer) GetAllDirectives() []directive.Directive {
	return m.directives
}

// GetLastRoutingDecision returns the most recent routing decision.
// Returns nil if no routing has been performed yet.
func (m *RouterMultiplexer) GetLastRoutingDecision() *router.RoutingDecision {
	return m.lastRoutingDecision
}
