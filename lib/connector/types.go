package connector

// ChatMessage represents a message in a chat conversation.
// This is a generic type used by connectors for LLM communication.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// ToolCallID is used for tool result messages
	ToolCallID string `json:"tool_call_id,omitempty"`
	// ToolCalls contains tool calls made by the LLM
	ToolCalls []ChatMessageToolCall `json:"tool_calls,omitempty"`
}

// ChatMessageToolCall represents a tool call in a chat message.
type ChatMessageToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function FunctionCallObj `json:"function"`
}

// FunctionCallObj represents a function call within a message.
type FunctionCallObj struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ResponseFormat specifies the format for structured output.
type ResponseFormat struct {
	Type       string      `json:"type"` // "json_schema", "json_object", "text"
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

func NewJSONSchemaResponseFormat(JSONSchema *JSONSchema) *ResponseFormat {
	return &ResponseFormat{
		Type:       "json_schema",
		JSONSchema: JSONSchema,
	}
}

// JSONSchema defines the schema for structured output.
type JSONSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Strict      bool           `json:"strict"`
	Schema      map[string]any `json:"schema"`
}
