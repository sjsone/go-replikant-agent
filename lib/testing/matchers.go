package testing

import (
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// AssertContextPartEqual asserts that two ContextParts are equal.
func AssertContextPartEqual(t *testing.T, expected, actual *agentic_context.ContextPart) bool {
	t.Helper()

	if expected == nil && actual == nil {
		return true
	}

	if expected == nil {
		t.Errorf("Expected nil context part, got non-nil")
		return false
	}

	if actual == nil {
		t.Errorf("Expected non-nil context part, got nil")
		return false
	}

	if expected.Raw != actual.Raw {
		t.Errorf("Expected Raw to be %q, got %q", expected.Raw, actual.Raw)
		return false
	}

	if expected.Source != actual.Source {
		t.Errorf("Expected Source to be %v, got %v", expected.Source, actual.Source)
		return false
	}

	if expected.ToolUse != actual.ToolUse {
		t.Errorf("Expected ToolUse to be %v, got %v", expected.ToolUse, actual.ToolUse)
		return false
	}

	if expected.Stop != actual.Stop {
		t.Errorf("Expected Stop to be %v, got %v", expected.Stop, actual.Stop)
		return false
	}

	if expected.Cancelled != actual.Cancelled {
		t.Errorf("Expected Cancelled to be %v, got %v", expected.Cancelled, actual.Cancelled)
		return false
	}

	if !assertToolCallsEqual(t, expected.ToolCalls, actual.ToolCalls) {
		return false
	}

	if !assertToolResultsEqual(t, expected.ToolResults, actual.ToolResults) {
		return false
	}

	return true
}

// assertToolCallsEqual asserts that two slices of FunctionCall are equal.
func assertToolCallsEqual(t *testing.T, expected, actual []tool.FunctionCall) bool {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected %d tool calls, got %d", len(expected), len(actual))
		return false
	}

	for i := range expected {
		if expected[i].ID != actual[i].ID {
			t.Errorf("Expected tool call %d ID to be %q, got %q", i, expected[i].ID, actual[i].ID)
			return false
		}

		if expected[i].Name != actual[i].Name {
			t.Errorf("Expected tool call %d Name to be %q, got %q", i, expected[i].Name, actual[i].Name)
			return false
		}

		if !assertMapsEqual(t, expected[i].Arguments, actual[i].Arguments) {
			return false
		}
	}

	return true
}

// AssertToolCallsEqual asserts that two slices of FunctionCall are equal.
func AssertToolCallsEqual(t *testing.T, expected, actual []tool.FunctionCall) bool {
	t.Helper()
	return assertToolCallsEqual(t, expected, actual)
}

// assertToolResultsEqual asserts that two slices of FunctionResult are equal.
func assertToolResultsEqual(t *testing.T, expected, actual []tool.FunctionResult) bool {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected %d tool results, got %d", len(expected), len(actual))
		return false
	}

	for i := range expected {
		if expected[i].ID != actual[i].ID {
			t.Errorf("Expected tool result %d ID to be %q, got %q", i, expected[i].ID, actual[i].ID)
			return false
		}

		if expected[i].Content != actual[i].Content {
			t.Errorf("Expected tool result %d Content to be %q, got %q", i, expected[i].Content, actual[i].Content)
			return false
		}

		if expected[i].Error != actual[i].Error {
			t.Errorf("Expected tool result %d Error to be %v, got %v", i, expected[i].Error, actual[i].Error)
			return false
		}
	}

	return true
}

// AssertToolResultsEqual asserts that two slices of FunctionResult are equal.
func AssertToolResultsEqual(t *testing.T, expected, actual []tool.FunctionResult) bool {
	t.Helper()
	return assertToolResultsEqual(t, expected, actual)
}

// assertMapsEqual asserts that two maps are equal.
func assertMapsEqual(t *testing.T, expected, actual map[string]any) bool {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected map to have %d keys, got %d", len(expected), len(actual))
		return false
	}

	for k, expectedVal := range expected {
		actualVal, ok := actual[k]
		if !ok {
			t.Errorf("Expected map to have key %q", k)
			return false
		}

		if expectedVal != actualVal {
			t.Errorf("Expected map key %q to be %v, got %v", k, expectedVal, actualVal)
			return false
		}
	}

	return true
}

// AssertMapsEqual asserts that two maps are equal.
func AssertMapsEqual(t *testing.T, expected, actual map[string]any) bool {
	t.Helper()
	return assertMapsEqual(t, expected, actual)
}

// AssertNotNil asserts that a value is not nil.
func AssertNotNil(t *testing.T, value any, msg string) bool {
	t.Helper()

	if value == nil {
		t.Errorf("%s: expected non-nil, got nil", msg)
		return false
	}

	return true
}

// AssertNil asserts that a value is nil.
func AssertNil(t *testing.T, value any, msg string) bool {
	t.Helper()

	if value != nil {
		t.Errorf("%s: expected nil, got %v", msg, value)
		return false
	}

	return true
}

// AssertEqual asserts that two values are equal.
func AssertEqual[T comparable](t *testing.T, expected, actual T, msg string) bool {
	t.Helper()

	if expected != actual {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
		return false
	}

	return true
}

// AssertSliceEqual asserts that two slices are equal.
func AssertSliceEqual[T comparable](t *testing.T, expected, actual []T, msg string) bool {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("%s: expected length %d, got %d", msg, len(expected), len(actual))
		return false
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("%s: at index %d, expected %v, got %v", msg, i, expected[i], actual[i])
			return false
		}
	}

	return true
}

// AssertSliceLength asserts that a slice has a specific length.
func AssertSliceLength[T any](t *testing.T, slice []T, expectedLen int, msg string) bool {
	t.Helper()

	if len(slice) != expectedLen {
		t.Errorf("%s: expected length %d, got %d", msg, expectedLen, len(slice))
		return false
	}

	return true
}

// AssertContains asserts that a slice contains a value.
func AssertContains[T comparable](t *testing.T, slice []T, value T, msg string) bool {
	t.Helper()

	for _, v := range slice {
		if v == value {
			return true
		}
	}

	t.Errorf("%s: expected slice to contain %v", msg, value)
	return false
}

// AssertNotEmpty asserts that a string is not empty.
func AssertNotEmpty(t *testing.T, s string, msg string) bool {
	t.Helper()

	if s == "" {
		t.Errorf("%s: expected non-empty string", msg)
		return false
	}

	return true
}

// AssertTrue asserts that a boolean is true.
func AssertTrue(t *testing.T, value bool, msg string) bool {
	t.Helper()

	if !value {
		t.Errorf("%s: expected true, got false", msg)
		return false
	}

	return true
}

// AssertFalse asserts that a boolean is false.
func AssertFalse(t *testing.T, value bool, msg string) bool {
	t.Helper()

	if value {
		t.Errorf("%s: expected false, got true", msg)
		return false
	}

	return true
}

// AssertNoError asserts that an error is nil.
func AssertNoError(t *testing.T, err error, msg string) bool {
	t.Helper()

	if err != nil {
		t.Errorf("%s: expected no error, got %v", msg, err)
		return false
	}

	return true
}

// AssertError asserts that an error is not nil.
func AssertError(t *testing.T, err error, msg string) bool {
	t.Helper()

	if err == nil {
		t.Errorf("%s: expected error, got nil", msg)
		return false
	}

	return true
}

// AssertSourceIsUser asserts that a source is UserSource.
func AssertSourceIsUser(t *testing.T, source agentic_context.Source) bool {
	t.Helper()

	if !source.IsUser() {
		t.Errorf("Expected source to be UserSource")
		return false
	}

	return true
}

// AssertSourceIsAgent asserts that a source is AgentSource.
func AssertSourceIsAgent(t *testing.T, source agentic_context.Source) bool {
	t.Helper()

	if !source.IsAgent() {
		t.Errorf("Expected source to be AgentSource")
		return false
	}

	return true
}

// AssertSourceIsSystem asserts that a source is SystemSource.
func AssertSourceIsSystem(t *testing.T, source agentic_context.Source) bool {
	t.Helper()

	if !source.IsSystem() {
		t.Errorf("Expected source to be SystemSource")
		return false
	}

	return true
}

// AssertSourceIsTool asserts that a source is ToolSource.
func AssertSourceIsTool(t *testing.T, source agentic_context.Source) bool {
	t.Helper()

	if !source.IsTool() {
		t.Errorf("Expected source to be ToolSource")
		return false
	}

	return true
}
