package router

// RoutingDecision represents a routing decision for directive selection.
type RoutingDecision struct {
	SelectedIDs []string `json:"selected_ids"`
	Reasoning   string   `json:"reasoning"`
	Confidence  float64  `json:"confidence"`
}

// RoutingDecisionParams defines the schema for routing decisions.
type RoutingDecisionParams struct {
	SelectedIDs []string `json:"selected_ids" jsonschema:"MUST contain ALL option names that should be activated. If multiple tools work together, include ALL their names here"`
	Reasoning   string   `json:"reasoning" jsonschema:"Explanation for why these directives were selected"`
}
