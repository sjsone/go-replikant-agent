package router

import "context"

type RoutingOption struct {
	Name      string
	Text      string
	ToolNames []string // Names of tools provided by this directive
}

// RoutingResult contains both the selected options and the routing decision.
type RoutingResult struct {
	SelectedOptions []*RoutingOption
	Decision        *RoutingDecision
}

type Router interface {
	SetDelegate(d Delegate)
	Route(ctx context.Context, userQuery string, allAvailableOptions []*RoutingOption) *RoutingResult
}
