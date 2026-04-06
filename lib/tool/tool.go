package tool

import (
	"context"
)

// Tool represents a callable function that the LLM can invoke.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"` // JSON Schema
}

// ToolCallable is the interface that tools must implement for execution.
type ToolCallable interface {
	Execute(ctx context.Context, args map[string]any) (string, error)
	GetName() string
}

// FunctionCall represents a parsed tool call from the LLM response.
type FunctionCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// FunctionResult represents the result of executing a tool.
type FunctionResult struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Error   bool   `json:"error,omitempty"`
}
