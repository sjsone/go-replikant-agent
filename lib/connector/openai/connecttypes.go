package openai

import (
	libconnector "github.com/sjsone/go-replikant-agent/lib/connector"
)

// chatMessage is the OpenAI wire format for messages.
type chatMessage struct {
	Role       string             `json:"role"`
	Content    string             `json:"content"`
	ToolCallID string             `json:"tool_call_id,omitempty"`
	ToolCalls  []chatMessageToolCall `json:"tool_calls,omitempty"`
}

// chatMessageToolCall represents a tool call in a chat message.
type chatMessageToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function functionCallObj `json:"function"`
}

// functionCallObj represents a function call within a message.
type functionCallObj struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatRequest represents the request body for OpenAI chat completion API.
type ChatRequest struct {
	Model          string           `json:"model"`
	Messages       []chatMessage    `json:"messages"`
	Stream         bool             `json:"stream,omitempty"`
	Tools          []ToolDefinition `json:"tools,omitempty"`
	ResponseFormat *ResponseFormat  `json:"response_format,omitempty"`
}

// Choice represents a single completion choice in the API response.
type Choice struct {
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatResponse represents the response from OpenAI chat completion API.
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// ErrorResponse represents an error response from the OpenAI API.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// StreamDelta represents a single delta in streaming response.
type StreamDelta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// StreamChoice represents a single choice in streaming response.
type StreamChoice struct {
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
	Index        int         `json:"index"`
}

// StreamResponse represents a single SSE chunk from the streaming API.
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// FunctionDefinition represents a function definition for OpenAI function calling.
type FunctionDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// ToolDefinition represents a tool in OpenAI format.
type ToolDefinition struct {
	Type     string             `json:"type"` // "function"
	Function FunctionDefinition `json:"function"`
}

// FunctionCall represents a function call in a response.
type FunctionCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"` // "function"
	Function functionCallObj `json:"function"`
}

// ToolCall represents a tool call in streaming response.
type ToolCall struct {
	Index    int             `json:"index"`
	ID       string          `json:"id,omitempty"`
	Type     string          `json:"type,omitempty"`
	Function functionCallObj `json:"function,omitempty"`
}

// ResponseFormat is an alias to the lib/connector.ResponseFormat type.
type ResponseFormat = libconnector.ResponseFormat

// JSONSchema is an alias to the lib/connector.JSONSchema type.
type JSONSchema = libconnector.JSONSchema
