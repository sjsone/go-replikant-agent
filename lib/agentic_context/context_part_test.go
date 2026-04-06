package agentic_context

import (
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

func TestNewSystemContextPart(t *testing.T) {
	raw := "System message"
	part := NewSystemContextPart(raw)

	if part == nil {
		t.Fatal("NewSystemContextPart returned nil")
	}

	if part.Raw != raw {
		t.Errorf("Expected Raw to be %q, got %q", raw, part.Raw)
	}

	if part.Source != SystemSource {
		t.Errorf("Expected Source to be SystemSource, got %v", part.Source)
	}
}

func TestNewAgentContextPart(t *testing.T) {
	raw := "Agent message"
	part := NewAgentContextPart(raw)

	if part == nil {
		t.Fatal("NewAgentContextPart returned nil")
	}

	if part.Raw != raw {
		t.Errorf("Expected Raw to be %q, got %q", raw, part.Raw)
	}

	if part.Source != AgentSource {
		t.Errorf("Expected Source to be AgentSource, got %v", part.Source)
	}
}

func TestNewUserContextPart(t *testing.T) {
	raw := "User message"
	part := NewUserContextPart(raw)

	if part == nil {
		t.Fatal("NewUserContextPart returned nil")
	}

	if part.Raw != raw {
		t.Errorf("Expected Raw to be %q, got %q", raw, part.Raw)
	}

	if part.Source != UserSource {
		t.Errorf("Expected Source to be UserSource, got %v", part.Source)
	}
}

func TestNewToolResultContextPart(t *testing.T) {
	results := []tool.FunctionResult{
		{
			ID:      "call-1",
			Content: "Result 1",
			Error:   false,
		},
		{
			ID:      "call-2",
			Content: "Result 2",
			Error:   false,
		},
	}

	part := NewToolResultContextPart(results)

	if part == nil {
		t.Fatal("NewToolResultContextPart returned nil")
	}

	if part.Source != ToolSource {
		t.Errorf("Expected Source to be ToolSource, got %v", part.Source)
	}

	if len(part.ToolResults) != 2 {
		t.Errorf("Expected 2 tool results, got %d", len(part.ToolResults))
	}

	if part.ToolResults[0].ID != "call-1" {
		t.Errorf("Expected first result ID to be 'call-1', got %q", part.ToolResults[0].ID)
	}

	if part.ToolResults[1].Content != "Result 2" {
		t.Errorf("Expected second result content to be 'Result 2', got %q", part.ToolResults[1].Content)
	}

	// Check that Raw is formatted correctly
	expectedRaw := "Tool call-1: Result 1\nTool call-2: Result 2\n"
	if part.Raw != expectedRaw {
		t.Errorf("Expected Raw to be %q, got %q", expectedRaw, part.Raw)
	}
}

func TestNewToolResultContextPart_Empty(t *testing.T) {
	results := []tool.FunctionResult{}
	part := NewToolResultContextPart(results)

	if part == nil {
		t.Fatal("NewToolResultContextPart returned nil")
	}

	if part.Source != ToolSource {
		t.Errorf("Expected Source to be ToolSource, got %v", part.Source)
	}

	if len(part.ToolResults) != 0 {
		t.Errorf("Expected 0 tool results, got %d", len(part.ToolResults))
	}

	if part.Raw != "" {
		t.Errorf("Expected Raw to be empty, got %q", part.Raw)
	}
}

func TestNewToolResultContextPart_WithError(t *testing.T) {
	results := []tool.FunctionResult{
		{
			ID:      "call-1",
			Content: "Error: something went wrong",
			Error:   true,
		},
	}

	part := NewToolResultContextPart(results)

	if !part.ToolResults[0].Error {
		t.Error("Expected tool result to have Error=true")
	}

	if part.ToolResults[0].Content != "Error: something went wrong" {
		t.Errorf("Expected error content, got %q", part.ToolResults[0].Content)
	}
}

func TestIsToolResult(t *testing.T) {
	tests := []struct {
		name     string
		part     *ContextPart
		expected bool
	}{
		{
			name:     "context part with tool results",
			part:     NewToolResultContextPart([]tool.FunctionResult{{ID: "call-1", Content: "result"}}),
			expected: true,
		},
		{
			name:     "context part with empty tool results",
			part:     NewToolResultContextPart([]tool.FunctionResult{}),
			expected: false,
		},
		{
			name:     "user context part",
			part:     NewUserContextPart("hello"),
			expected: false,
		},
		{
			name:     "agent context part",
			part:     NewAgentContextPart("hello"),
			expected: false,
		},
		{
			name:     "system context part",
			part:     NewSystemContextPart("hello"),
			expected: false,
		},
		{
			name: "context part with tool calls but no results",
			part: &ContextPart{
				Raw:       "tool use",
				Source:    AgentSource,
				ToolCalls: []tool.FunctionCall{{ID: "call-1", Name: "test"}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.part.IsToolResult() != tt.expected {
				t.Errorf("IsToolResult() = %v, want %v", tt.part.IsToolResult(), tt.expected)
			}
		})
	}
}

func TestContextPartFields(t *testing.T) {
	part := &ContextPart{
		Raw:         "test content",
		Source:      UserSource,
		ToolUse:     true,
		Stop:        true,
		Cancelled:   true,
		ToolCalls:   []tool.FunctionCall{{ID: "call-1", Name: "tool1"}},
		ToolResults: []tool.FunctionResult{{ID: "call-1", Content: "result"}},
	}

	tests := []struct {
		name     string
		field    any
		expected any
	}{
		{"Raw", part.Raw, "test content"},
		{"Source", part.Source, UserSource},
		{"ToolUse", part.ToolUse, true},
		{"Stop", part.Stop, true},
		{"Cancelled", part.Cancelled, true},
		{"ToolCalls", len(part.ToolCalls), 1},
		{"ToolResults", len(part.ToolResults), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.field)
			}
		})
	}
}

func TestContextPart_DefaultValues(t *testing.T) {
	part := &ContextPart{
		Raw:    "test",
		Source: UserSource,
	}

	if part.ToolUse {
		t.Error("Expected ToolUse to default to false")
	}

	if part.Stop {
		t.Error("Expected Stop to default to false")
	}

	if part.Cancelled {
		t.Error("Expected Cancelled to default to false")
	}

	if len(part.ToolCalls) != 0 {
		t.Error("Expected ToolCalls to default to empty slice")
	}

	if len(part.ToolResults) != 0 {
		t.Error("Expected ToolResults to default to empty slice")
	}
}

func TestContextPart_ToolResultsFormatting(t *testing.T) {
	tests := []struct {
		name     string
		results  []tool.FunctionResult
		expected string
	}{
		{
			name: "single result",
			results: []tool.FunctionResult{
				{ID: "tool-1", Content: "Success"},
			},
			expected: "Tool tool-1: Success\n",
		},
		{
			name: "multiple results",
			results: []tool.FunctionResult{
				{ID: "tool-1", Content: "Result 1"},
				{ID: "tool-2", Content: "Result 2"},
				{ID: "tool-3", Content: "Result 3"},
			},
			expected: "Tool tool-1: Result 1\nTool tool-2: Result 2\nTool tool-3: Result 3\n",
		},
		{
			name:     "empty results",
			results:  []tool.FunctionResult{},
			expected: "",
		},
		{
			name: "result with error",
			results: []tool.FunctionResult{
				{ID: "tool-1", Content: "Error: failed", Error: true},
			},
			expected: "Tool tool-1: Error: failed\n",
		},
		{
			name: "result with multiline content",
			results: []tool.FunctionResult{
				{ID: "tool-1", Content: "Line 1\nLine 2\nLine 3"},
			},
			expected: "Tool tool-1: Line 1\nLine 2\nLine 3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part := NewToolResultContextPart(tt.results)
			if part.Raw != tt.expected {
				t.Errorf("Expected Raw to be %q, got %q", tt.expected, part.Raw)
			}
		})
	}
}
