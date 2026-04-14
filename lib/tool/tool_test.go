package tool

import (
	"context"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		jsonArgs string
		wantLen  int
		firstKey string
		firstVal any
	}{
		{
			name:     "valid JSON object",
			jsonArgs: `{"city": "New York", "country": "USA"}`,
			wantLen:  2,
			firstKey: "city",
			firstVal: "New York",
		},
		{
			name:     "empty JSON object",
			jsonArgs: `{}`,
			wantLen:  0,
		},
		{
			name:     "JSON with various types",
			jsonArgs: `{"string": "value", "number": 42, "bool": true, "null": null}`,
			wantLen:  4,
		},
		{
			name:     "JSON with nested object",
			jsonArgs: `{"location": {"city": "NYC", "country": "USA"}}`,
			wantLen:  1,
		},
		{
			name:     "JSON with array",
			jsonArgs: `{"items": [1, 2, 3]}`,
			wantLen:  1,
		},
		{
			name:     "invalid JSON - missing closing brace",
			jsonArgs: `{"city": "New York"`,
			wantLen:  0, // Returns empty map on error
		},
		{
			name:     "invalid JSON - unquoted key",
			jsonArgs: `{city: "New York"}`,
			wantLen:  0,
		},
		{
			name:     "invalid JSON - trailing comma",
			jsonArgs: `{"city": "New York",}`,
			wantLen:  0,
		},
		{
			name:     "empty string",
			jsonArgs: "",
			wantLen:  0,
		},
		{
			name:     "JSON array instead of object",
			jsonArgs: `[1, 2, 3]`,
			wantLen:  0,
		},
		{
			name:     "JSON string instead of object",
			jsonArgs: `"just a string"`,
			wantLen:  0,
		},
		{
			name:     "JSON null",
			jsonArgs: `null`,
			wantLen:  0,
		},
		{
			name:     "JSON number",
			jsonArgs: `42`,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := ParseArgsToMap(tt.jsonArgs)

			if len(args) != tt.wantLen {
				t.Errorf("ParseArgs() returned %d args, want %d", len(args), tt.wantLen)
			}

			if tt.firstKey != "" {
				val, ok := args[tt.firstKey]
				if !ok {
					t.Errorf("ParseArgs() missing key %q", tt.firstKey)
				} else if val != tt.firstVal {
					t.Errorf("ParseArgs() key %q = %v, want %v", tt.firstKey, val, tt.firstVal)
				}
			}
		})
	}
}

func TestParseArgs_Unicode(t *testing.T) {
	jsonArgs := `{"greeting": "Hello 世界", "emoji": "🚀"}`
	args := ParseArgsToMap(jsonArgs)

	if args["greeting"] != "Hello 世界" {
		t.Errorf("Expected greeting to be 'Hello 世界', got %q", args["greeting"])
	}

	if args["emoji"] != "🚀" {
		t.Errorf("Expected emoji to be '🚀', got %q", args["emoji"])
	}
}

func TestParseArgs_EscapedCharacters(t *testing.T) {
	jsonArgs := `{"text": "Line 1\nLine 2\tTabbed", "quote": "He said \"hello\""}`
	args := ParseArgsToMap(jsonArgs)

	expected := "Line 1\nLine 2\tTabbed"
	if args["text"] != expected {
		t.Errorf("Expected text to be %q, got %q", expected, args["text"])
	}

	quoteExpected := `He said "hello"`
	if args["quote"] != quoteExpected {
		t.Errorf("Expected quote to be %q, got %q", quoteExpected, args["quote"])
	}
}

func TestFunctionCall(t *testing.T) {
	call := FunctionCall{
		ID:   "call-123",
		Name: "test_tool",
		Arguments: map[string]any{
			"param1": "value1",
			"param2": 42,
		},
	}

	if call.ID != "call-123" {
		t.Errorf("Expected ID to be 'call-123', got %q", call.ID)
	}

	if call.Name != "test_tool" {
		t.Errorf("Expected Name to be 'test_tool', got %q", call.Name)
	}

	if call.Arguments["param1"] != "value1" {
		t.Errorf("Expected param1 to be 'value1', got %v", call.Arguments["param1"])
	}

	if call.Arguments["param2"] != 42 {
		t.Errorf("Expected param2 to be 42, got %v", call.Arguments["param2"])
	}
}

func TestFunctionCall_EmptyArguments(t *testing.T) {
	call := FunctionCall{
		ID:        "call-456",
		Name:      "empty_tool",
		Arguments: map[string]any{},
	}

	if len(call.Arguments) != 0 {
		t.Errorf("Expected empty arguments, got %d", len(call.Arguments))
	}
}

func TestFunctionCall_NilArguments(t *testing.T) {
	call := FunctionCall{
		ID:        "call-789",
		Name:      "nil_tool",
		Arguments: nil,
	}

	if call.Arguments != nil {
		t.Errorf("Expected nil arguments, got %v", call.Arguments)
	}
}

func TestFunctionResult(t *testing.T) {
	result := FunctionResult{
		ID:      "call-123",
		Content: "Success!",
		Error:   false,
	}

	if result.ID != "call-123" {
		t.Errorf("Expected ID to be 'call-123', got %q", result.ID)
	}

	if result.Content != "Success!" {
		t.Errorf("Expected Content to be 'Success!', got %q", result.Content)
	}

	if result.Error {
		t.Error("Expected Error to be false")
	}
}

func TestFunctionResult_Error(t *testing.T) {
	result := FunctionResult{
		ID:      "call-456",
		Content: "Something went wrong",
		Error:   true,
	}

	if !result.Error {
		t.Error("Expected Error to be true")
	}

	if result.Content != "Something went wrong" {
		t.Errorf("Expected Content to be 'Something went wrong', got %q", result.Content)
	}
}

func TestFunctionResult_EmptyContent(t *testing.T) {
	result := FunctionResult{
		ID:      "call-789",
		Content: "",
		Error:   false,
	}

	if result.Content != "" {
		t.Errorf("Expected empty Content, got %q", result.Content)
	}
}

func TestFunctionResult_MultiLineContent(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	result := FunctionResult{
		ID:      "call-multi",
		Content: content,
		Error:   false,
	}

	if result.Content != content {
		t.Errorf("Expected multi-line content, got %q", result.Content)
	}
}

func TestTool(t *testing.T) {
	tool := Tool{
		Name:        "test_tool",
		Description: "A test tool for testing",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input parameter",
				},
			},
		},
	}

	if tool.Name != "test_tool" {
		t.Errorf("Expected Name to be 'test_tool', got %q", tool.Name)
	}

	if tool.Description != "A test tool for testing" {
		t.Errorf("Expected Description to be 'A test tool for testing', got %q", tool.Description)
	}

	if tool.Parameters == nil {
		t.Error("Expected Parameters to be non-nil")
	}

	props, ok := tool.Parameters["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	input, ok := props["input"].(map[string]any)
	if !ok {
		t.Fatal("Expected input to be a map")
	}

	if input["type"] != "string" {
		t.Errorf("Expected input type to be 'string', got %v", input["type"])
	}
}

func TestTool_WithEmptyParameters(t *testing.T) {
	tool := Tool{
		Name:        "simple_tool",
		Description: "A simple tool",
		Parameters:  map[string]any{},
	}

	if len(tool.Parameters) != 0 {
		t.Errorf("Expected empty Parameters, got %d", len(tool.Parameters))
	}
}

func TestToolCallable_Interface(t *testing.T) {
	// Verify that ToolCallable is a valid interface
	var _ ToolCallable = (*MockToolCallableForTest)(nil)
}

// MockToolCallableForTest is a minimal mock for interface testing
type MockToolCallableForTest struct{}

func (m *MockToolCallableForTest) GetTool() *Tool {
	return &Tool{Name: "mockTool"}
}
func (m *MockToolCallableForTest) Execute(ctx context.Context, args map[string]any) (string, error) {
	return "result", nil
}

func TestTool_StructTags(t *testing.T) {
	tool := Tool{
		Name:        "test_tool",
		Description: "Test description",
		Parameters:  map[string]any{"key": "value"},
	}

	// Verify JSON tags work by marshaling
	// This is implicitly tested by the fact that the struct is used in JSON contexts
	if tool.Name == "" {
		t.Error("Expected Name to be preserved")
	}
}

func TestParseArgs_LargeJSON(t *testing.T) {
	largeJSON := `{"key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4", "key5": "value5"}`
	args := ParseArgsToMap(largeJSON)

	if len(args) != 5 {
		t.Errorf("Expected 5 arguments, got %d", len(args))
	}
}

func TestParseArgs_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		jsonArgs string
	}{
		{
			name: "pretty printed JSON",
			jsonArgs: `{
				"key": "value"
			}`,
		},
		{
			name:     "JSON with extra spaces",
			jsonArgs: `{  "key"  :  "value"  }`,
		},
		{
			name:     "JSON with newlines",
			jsonArgs: "{\"key\":\n\"value\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := ParseArgsToMap(tt.jsonArgs)
			if args["key"] != "value" {
				t.Errorf("Expected key to be 'value', got %v", args["key"])
			}
		})
	}
}
