package string

import (
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

func TestNewStringPromptBuilder(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
	}{
		{
			name:   "simple prompt",
			prompt: "You are a helpful assistant",
		},
		{
			name:   "empty prompt",
			prompt: "",
		},
		{
			name:   "multiline prompt",
			prompt: "You are a helpful assistant.\n\nYou have access to tools.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewStringPromptBuilder(tt.prompt)

			if pb == nil {
				t.Fatal("NewStringPromptBuilder returned nil")
			}

			if pb.raw != tt.prompt {
				t.Errorf("raw field = %q, want %q", pb.raw, tt.prompt)
			}
		})
	}
}

func TestStringPromptBuilder_Build(t *testing.T) {
	tests := []struct {
		name       string
		basePrompt string
		directives []directive.Directive
		expected   string
	}{
		{
			name:       "no directives",
			basePrompt: "Base prompt",
			directives: []directive.Directive{},
			expected:   "Base prompt",
		},
		{
			name:       "nil directives",
			basePrompt: "Base prompt",
			directives: nil,
			expected:   "Base prompt",
		},
		{
			name:       "single directive",
			basePrompt: "Base prompt",
			directives: []directive.Directive{
				newTestDirective("Directive 1"),
			},
			expected: "Base prompt\nDirective 1",
		},
		{
			name:       "multiple directives",
			basePrompt: "Base prompt",
			directives: []directive.Directive{
				newTestDirective("Directive 1"),
				newTestDirective("Directive 2"),
				newTestDirective("Directive 3"),
			},
			expected: "Base prompt\nDirective 1\nDirective 2\nDirective 3",
		},
		{
			name:       "directives with newlines",
			basePrompt: "Base",
			directives: []directive.Directive{
				newTestDirective("Line 1\nLine 2"),
				newTestDirective("Line 3\nLine 4"),
			},
			expected: "Base\nLine 1\nLine 2\nLine 3\nLine 4",
		},
		{
			name:       "empty base prompt",
			basePrompt: "",
			directives: []directive.Directive{
				newTestDirective("Directive 1"),
			},
			expected: "\nDirective 1",
		},
		{
			name:       "empty directive prompt",
			basePrompt: "Base",
			directives: []directive.Directive{
				newTestDirective(""),
			},
			expected: "Base\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewStringPromptBuilder(tt.basePrompt)
			result := pb.Build(tt.directives)

			if result.Raw != tt.expected {
				t.Errorf("Build() Raw = %q, want %q", result.Raw, tt.expected)
			}
		})
	}
}

func TestStringPromptBuilder_BuildWithNewlines(t *testing.T) {
	basePrompt := "You are a helpful assistant."

	directives := []directive.Directive{
		newTestDirective("You have access to the following tools:\n- tool1\n- tool2"),
		newTestDirective("\nAlways use tools when appropriate."),
		newTestDirective("End of instructions."),
	}

	pb := NewStringPromptBuilder(basePrompt)
	result := pb.Build(directives)

	expected := "You are a helpful assistant.\n" +
		"You have access to the following tools:\n- tool1\n- tool2\n" +
		"\nAlways use tools when appropriate.\n" +
		"End of instructions."

	if result.Raw != expected {
		t.Errorf("Build() Raw =\n%q\nwant\n%q", result.Raw, expected)
	}
}

func TestStringPromptBuilder_BuildIsDeterministic(t *testing.T) {
	basePrompt := "Base"
	directives := []directive.Directive{
		newTestDirective("Directive 1"),
		newTestDirective("Directive 2"),
	}

	pb := NewStringPromptBuilder(basePrompt)
	result1 := pb.Build(directives)
	result2 := pb.Build(directives)

	if result1.Raw != result2.Raw {
		t.Error("Build() is not deterministic")
	}
}

func TestStringPromptBuilder_ModifyOriginalPrompt(t *testing.T) {
	originalPrompt := "Original"
	pb := NewStringPromptBuilder(originalPrompt)

	// Strings are immutable in Go, so reassigning the variable
	// doesn't affect the builder's internal state
	originalPrompt = "Modified"

	// The builder should still have the original value
	directives := []directive.Directive{}
	result := pb.Build(directives)

	if result.Raw != "Original" {
		t.Errorf("Build() Raw = %q, want 'Original'", result.Raw)
	}
}

func TestStringPromptBuilder_EmptyBaseWithEmptyDirectives(t *testing.T) {
	pb := NewStringPromptBuilder("")
	result := pb.Build([]directive.Directive{})

	if result.Raw != "" {
		t.Errorf("Build() Raw = %q, want empty string", result.Raw)
	}
}

func TestStringPromptBuilder_DirectiveOrderPreserved(t *testing.T) {
	basePrompt := "Base"

	directives := []directive.Directive{
		newTestDirective("First"),
		newTestDirective("Second"),
		newTestDirective("Third"),
	}

	pb := NewStringPromptBuilder(basePrompt)
	result := pb.Build(directives)

	// Check that directives are in the correct order
	expected := "Base\nFirst\nSecond\nThird"
	if result.Raw != expected {
		t.Errorf("Build() Raw = %q, want %q", result.Raw, expected)
	}
}

func TestStringPromptBuilder_ReturnsPromptStruct(t *testing.T) {
	pb := NewStringPromptBuilder("Base")
	result := pb.Build([]directive.Directive{})

	// Check that result is a Prompt struct
	if result.Raw != "Base" {
		t.Errorf("Expected Raw to be 'Base', got %q", result.Raw)
	}
}

// Helper function to create test directives
func newTestDirective(promptText string) *directive.StaticDirective {
	p := &prompt.Prompt{Raw: promptText}
	return directive.NewStaticDirective("test", p, []tool.ToolCallable{})
}
