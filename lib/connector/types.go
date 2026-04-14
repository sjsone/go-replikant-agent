package connector

// ResponseFormat specifies the format for structured output.
type ResponseFormat struct {
	Type       string      `json:"type"` // "json_schema", "json_object", "text"
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

// JSONSchema defines the schema for structured output.
type JSONSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Strict      bool           `json:"strict"`
	Schema      map[string]any `json:"schema"`
}
