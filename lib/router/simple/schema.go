package simple

import "github.com/sjsone/go-replikant-agent/lib/connector"

// BuildRoutingDecisionSchema builds a JSON schema for routing decisions
// with an enum constraint on selected_ids items, restricting them to the
// provided option names.
func BuildRoutingDecisionSchema(optionNames []string) *connector.JSONSchema {
	// Convert string slice to []any for the enum field
	enumValues := make([]any, len(optionNames))
	for i, name := range optionNames {
		enumValues[i] = name
	}

	return &connector.JSONSchema{
		Name:        "routing decision",
		Strict:      true,
		Description: "",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"selected_ids": map[string]any{
					"type":        "array",
					"items":       map[string]any{"type": "string", "enum": enumValues},
					"description": "MUST contain ALL option names that should be activated",
				},
				"reasoning": map[string]any{
					"type":        "string",
					"description": "Explanation for why these options were selected",
				},
			},
			"required":             []string{"selected_ids", "reasoning"},
			"additionalProperties": false,
		},
	}
}
