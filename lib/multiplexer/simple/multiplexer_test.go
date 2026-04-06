package multiplexer

import (
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

func TestNewSimpleMultiplexer(t *testing.T) {
	tests := []struct {
		name       string
		directives []directive.Directive
	}{
		{
			name:       "with directives",
			directives: []directive.Directive{newTestDirective("d1"), newTestDirective("d2")},
		},
		{
			name:       "with nil directives",
			directives: nil,
		},
		{
			name:       "with empty directives",
			directives: []directive.Directive{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewSimpleMultiplexer(tt.directives)

			if m == nil {
				t.Fatal("NewSimpleMultiplexer returned nil")
			}

			if tt.directives == nil {
				if m.directives != nil {
					t.Errorf("directives field = %v, want nil", m.directives)
				}
			} else if len(m.directives) != len(tt.directives) {
				t.Errorf("directives length = %d, want %d", len(m.directives), len(tt.directives))
			}
		})
	}
}

func TestSimpleMultiplexer_GetActiveDirectivesForContext(t *testing.T) {
	tests := []struct {
		name       string
		directives []directive.Directive
		context    *agentic_context.AgentContext
		expected   int
	}{
		{
			name:       "returns all directives",
			directives: []directive.Directive{newTestDirective("d1"), newTestDirective("d2"), newTestDirective("d3")},
			context:    agentic_context.NewAgentContext(),
			expected:   3,
		},
		{
			name:       "returns empty directives",
			directives: []directive.Directive{},
			context:    agentic_context.NewAgentContext(),
			expected:   0,
		},
		{
			name:       "returns nil directives",
			directives: nil,
			context:    agentic_context.NewAgentContext(),
			expected:   0,
		},
		{
			name:       "context with user message",
			directives: []directive.Directive{newTestDirective("d1"), newTestDirective("d2")},
			context:    createContextWithParts(agentic_context.NewUserContextPart("Hello")),
			expected:   2,
		},
		{
			name:       "context with agent message",
			directives: []directive.Directive{newTestDirective("d1")},
			context:    createContextWithParts(agentic_context.NewAgentContextPart("Hi")),
			expected:   1,
		},
		{
			name:       "context with multiple parts",
			directives: []directive.Directive{newTestDirective("d1"), newTestDirective("d2")},
			context:    createContextWithParts(agentic_context.NewUserContextPart("Hello"), agentic_context.NewAgentContextPart("Hi")),
			expected:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewSimpleMultiplexer(tt.directives)
			result := m.GetActiveDirectivesForContext(*tt.context)

			if len(result) != tt.expected {
				t.Errorf("GetActiveDirectivesForContext() returned %d directives, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestSimpleMultiplexer_GetAllDirectives(t *testing.T) {
	tests := []struct {
		name       string
		directives []directive.Directive
		expected   int
	}{
		{
			name:       "multiple directives",
			directives: []directive.Directive{newTestDirective("d1"), newTestDirective("d2"), newTestDirective("d3")},
			expected:   3,
		},
		{
			name:       "single directive",
			directives: []directive.Directive{newTestDirective("d1")},
			expected:   1,
		},
		{
			name:       "empty directives",
			directives: []directive.Directive{},
			expected:   0,
		},
		{
			name:       "nil directives",
			directives: nil,
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewSimpleMultiplexer(tt.directives)
			result := m.GetAllDirectives()

			if len(result) != tt.expected {
				t.Errorf("GetAllDirectives() returned %d directives, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestSimpleMultiplexer_EmptyDirectives(t *testing.T) {
	m := NewSimpleMultiplexer([]directive.Directive{})
	ctx := agentic_context.NewAgentContext()

	result := m.GetActiveDirectivesForContext(*ctx)

	if len(result) != 0 {
		t.Errorf("GetActiveDirectivesForContext() returned %d directives, want 0", len(result))
	}

	allDirectives := m.GetAllDirectives()

	if len(allDirectives) != 0 {
		t.Errorf("GetAllDirectives() returned %d directives, want 0", len(allDirectives))
	}
}

func TestSimpleMultiplexer_ContextIndependence(t *testing.T) {
	directives := []directive.Directive{newTestDirective("d1"), newTestDirective("d2")}
	m := NewSimpleMultiplexer(directives)

	ctx1 := createContextWithParts(agentic_context.NewUserContextPart("Hello"))
	ctx2 := createContextWithParts(agentic_context.NewAgentContextPart("Hi"))
	ctx3 := agentic_context.NewAgentContext()

	result1 := m.GetActiveDirectivesForContext(*ctx1)
	result2 := m.GetActiveDirectivesForContext(*ctx2)
	result3 := m.GetActiveDirectivesForContext(*ctx3)

	// All should return the same directives regardless of context
	if len(result1) != len(result2) || len(result2) != len(result3) {
		t.Error("GetActiveDirectivesForContext() returned different lengths for different contexts")
	}

	if len(result1) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(result1))
	}
}

func TestSimpleMultiplexer_DirectiveOrderPreserved(t *testing.T) {
	d1 := newTestDirective("directive-1")
	d2 := newTestDirective("directive-2")
	d3 := newTestDirective("directive-3")

	directives := []directive.Directive{d1, d2, d3}
	m := NewSimpleMultiplexer(directives)

	ctx := agentic_context.NewAgentContext()
	result := m.GetActiveDirectivesForContext(*ctx)

	if len(result) != 3 {
		t.Fatalf("Expected 3 directives, got %d", len(result))
	}

	if result[0].GetName() != "directive-1" {
		t.Errorf("First directive name = %q, want 'directive-1'", result[0].GetName())
	}

	if result[1].GetName() != "directive-2" {
		t.Errorf("Second directive name = %q, want 'directive-2'", result[1].GetName())
	}

	if result[2].GetName() != "directive-3" {
		t.Errorf("Third directive name = %q, want 'directive-3'", result[2].GetName())
	}
}

// Helper functions

func newTestDirective(name string) *directive.StaticDirective {
	p := &prompt.Prompt{Raw: "Test directive: " + name}
	return directive.NewStaticDirective(name, p, []tool.ToolCallable{})
}

func createContextWithParts(parts ...*agentic_context.ContextPart) *agentic_context.AgentContext {
	ctx := agentic_context.NewAgentContext()
	for _, part := range parts {
		ctx.AddPart(part)
	}
	return ctx
}
