package multiplexer

import (
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/router"
	libtesting "github.com/sjsone/go-replikant-agent/lib/testing"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

func TestNewRouterMultiplexer(t *testing.T) {
	directives := []directive.Directive{newRouterTestDirective("d1"), newRouterTestDirective("d2")}
	mockRouter := libtesting.NewMockRouter()

	m := NewRouterMultiplexer(directives, mockRouter)

	if m == nil {
		t.Fatal("NewRouterMultiplexer returned nil")
	}

	if len(m.GetAllDirectives()) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(m.GetAllDirectives()))
	}

	if m.GetLastRoutingDecision() != nil {
		t.Error("Expected lastRoutingDecision to be nil initially")
	}
}

func TestRouterMultiplexer_GetActiveDirectivesForContext_NoUserMessage(t *testing.T) {
	directives := []directive.Directive{
		newRouterTestDirective("d1"),
		newRouterTestDirective("d2"),
		newRouterTestDirective("d3"),
	}
	mockRouter := libtesting.NewMockRouter()
	m := NewRouterMultiplexer(directives, mockRouter)

	// Context with no user message (agent message only)
	ctx := createRouterContextWithParts(agentic_context.NewAgentContextPart("Agent response"))
	result := m.GetActiveDirectivesForContext(*ctx)

	// Should return all directives when no user message
	if len(result) != 3 {
		t.Errorf("Expected 3 directives, got %d", len(result))
	}
}

func TestRouterMultiplexer_GetActiveDirectivesForContext_WithRouting(t *testing.T) {
	directives := []directive.Directive{
		newRouterTestDirective("d1"),
		newRouterTestDirective("d2"),
	}
	mockRouter := libtesting.NewMockRouter()
	m := NewRouterMultiplexer(directives, mockRouter)

	// Set up routing decision to select only directive d1
	mockRouter.SetRoutingDecision(&router.RoutingDecision{
		SelectedIDs: []string{"d1"},
		Reasoning:   "Only d1 is needed",
		Confidence:  0.9,
	})

	ctx := createRouterContextWithParts(agentic_context.NewUserContextPart("Hello"))
	result := m.GetActiveDirectivesForContext(*ctx)

	if len(result) != 1 {
		t.Errorf("Expected 1 directive, got %d", len(result))
	}

	if result[0].GetName() != "d1" {
		t.Errorf("Expected d1, got %s", result[0].GetName())
	}

	// Check that routing decision was stored
	decision := m.GetLastRoutingDecision()
	if decision == nil {
		t.Fatal("Expected routing decision to be stored")
	}

	if len(decision.SelectedIDs) != 1 || decision.SelectedIDs[0] != "d1" {
		t.Errorf("Routing decision not stored correctly: %v", decision.SelectedIDs)
	}
}

func TestRouterMultiplexer_ExtractUserMessage(t *testing.T) {
	tests := []struct {
		name     string
		parts    []*agentic_context.ContextPart
		expected string
	}{
		{
			name:     "single user message",
			parts:    []*agentic_context.ContextPart{agentic_context.NewUserContextPart("Hello")},
			expected: "Hello",
		},
		{
			name: "user then agent",
			parts: []*agentic_context.ContextPart{
				agentic_context.NewUserContextPart("Hello"),
				agentic_context.NewAgentContextPart("Hi"),
			},
			expected: "Hello",
		},
		{
			name: "agent then user",
			parts: []*agentic_context.ContextPart{
				agentic_context.NewAgentContextPart("Hi"),
				agentic_context.NewUserContextPart("Hello"),
			},
			expected: "Hello",
		},
		{
			name:     "no user message",
			parts:    []*agentic_context.ContextPart{agentic_context.NewAgentContextPart("Agent only")},
			expected: "",
		},
		{
			name: "multiple user messages",
			parts: []*agentic_context.ContextPart{
				agentic_context.NewUserContextPart("First"),
				agentic_context.NewAgentContextPart("Agent"),
				agentic_context.NewUserContextPart("Second"),
			},
			expected: "Second",
		},
		{
			name:     "empty context",
			parts:    []*agentic_context.ContextPart{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directives := []directive.Directive{newRouterTestDirective("d1")}
			mockRouter := libtesting.NewMockRouter()
			m := NewRouterMultiplexer(directives, mockRouter)

			ctx := createRouterContextWithParts(tt.parts...)
			result := m.extractUserMessage(*ctx)

			if result != tt.expected {
				t.Errorf("extractUserMessage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRouterMultiplexer_FilterDirectives(t *testing.T) {
	d1 := newRouterTestDirective("d1")
	d2 := newRouterTestDirective("d2")
	d3 := newRouterTestDirective("d3")
	directives := []directive.Directive{d1, d2, d3}

	tests := []struct {
		name     string
		ids      []string
		expected int
	}{
		{
			name:     "all directives",
			ids:      []string{"d1", "d2", "d3"},
			expected: 3,
		},
		{
			name:     "single directive",
			ids:      []string{"d2"},
			expected: 1,
		},
		{
			name:     "empty IDs",
			ids:      []string{},
			expected: 0,
		},
		{
			name:     "non-sequential IDs",
			ids:      []string{"d1", "d3"},
			expected: 2,
		},
		{
			name:     "unknown IDs",
			ids:      []string{"unknown", "d1", "d2", "nonexistent"},
			expected: 2, // Only d1 and d2 are valid
		},
		{
			name:     "all unknown IDs",
			ids:      []string{"foo", "bar", "baz"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRouter := libtesting.NewMockRouter()
			m := NewRouterMultiplexer(directives, mockRouter)

			// Set up routing decision with the test IDs
			mockRouter.SetRoutingDecision(&router.RoutingDecision{
				SelectedIDs: tt.ids,
				Reasoning:   "Test filter",
			})

			ctx := createRouterContextWithParts(agentic_context.NewUserContextPart("Test"))
			result := m.GetActiveDirectivesForContext(*ctx)

			if len(result) != tt.expected {
				t.Errorf("GetActiveDirectivesForContext() returned %d directives, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestRouterMultiplexer_GetLastRoutingDecision(t *testing.T) {
	directives := []directive.Directive{newRouterTestDirective("d1")}
	mockRouter := libtesting.NewMockRouter()
	m := NewRouterMultiplexer(directives, mockRouter)

	// Initially should be nil
	decision := m.GetLastRoutingDecision()
	if decision != nil {
		t.Error("Expected initial routing decision to be nil")
	}

	// Set a decision
	expectedDecision := &router.RoutingDecision{
		SelectedIDs: []string{"d1"},
		Reasoning:   "Test",
		Confidence:  0.8,
	}
	mockRouter.SetRoutingDecision(expectedDecision)

	ctx := createRouterContextWithParts(agentic_context.NewUserContextPart("Test"))
	m.GetActiveDirectivesForContext(*ctx)

	// Should now return the decision
	decision = m.GetLastRoutingDecision()
	if decision == nil {
		t.Fatal("Expected routing decision to be stored")
	}

	if len(decision.SelectedIDs) != 1 || decision.SelectedIDs[0] != "d1" {
		t.Errorf("Routing decision not correct: %v", decision.SelectedIDs)
	}
}

func TestRouterMultiplexer_GetAllDirectives(t *testing.T) {
	directives := []directive.Directive{
		newRouterTestDirective("d1"),
		newRouterTestDirective("d2"),
		newRouterTestDirective("d3"),
	}
	mockRouter := libtesting.NewMockRouter()
	m := NewRouterMultiplexer(directives, mockRouter)

	result := m.GetAllDirectives()

	if len(result) != 3 {
		t.Errorf("GetAllDirectives() returned %d directives, want 3", len(result))
	}
}

func TestRouterMultiplexer_RoutingError(t *testing.T) {
	directives := []directive.Directive{newRouterTestDirective("d1"), newRouterTestDirective("d2")}
	mockRouter := libtesting.NewMockRouter()
	m := NewRouterMultiplexer(directives, mockRouter)

	// Set routing error
	// Note: In the new Router interface, errors are handled differently
	// We simulate errors by setting an empty routing result
	mockRouter.SetRoutingDecision(&router.RoutingDecision{
		SelectedIDs: []string{},
		Reasoning:   "Error simulation",
	})

	ctx := createRouterContextWithParts(agentic_context.NewUserContextPart("Test"))
	result := m.GetActiveDirectivesForContext(*ctx)

	// Should fallback to all directives
	if len(result) != 2 {
		t.Errorf("Expected 2 directives on error, got %d", len(result))
	}
}

func TestRouterMultiplexer_RoutingWithMultipleSelectedIDs(t *testing.T) {
	d1 := newRouterTestDirective("d1")
	d2 := newRouterTestDirective("d2")
	d3 := newRouterTestDirective("d3")
	directives := []directive.Directive{d1, d2, d3}

	mockRouter := libtesting.NewMockRouter()
	m := NewRouterMultiplexer(directives, mockRouter)

	// Set routing decision to select multiple directives
	mockRouter.SetRoutingDecision(&router.RoutingDecision{
		SelectedIDs: []string{"d1", "d3"},
		Reasoning:   "Both d1 and d3 are needed",
		Confidence:  0.95,
	})

	ctx := createRouterContextWithParts(agentic_context.NewUserContextPart("Test"))
	result := m.GetActiveDirectivesForContext(*ctx)

	if len(result) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(result))
	}

	if result[0].GetName() != "d1" || result[1].GetName() != "d3" {
		t.Errorf("Expected d1 and d3, got %s and %s", result[0].GetName(), result[1].GetName())
	}
}

// Helper functions

func newRouterTestDirective(name string) *directive.StaticDirective {
	p := &prompt.Prompt{Raw: "Test directive: " + name}
	return directive.NewStaticDirective(name, p, []tool.ToolCallable{})
}

func createRouterContextWithParts(parts ...*agentic_context.ContextPart) *agentic_context.AgentContext {
	ctx := agentic_context.NewAgentContext()
	for _, part := range parts {
		ctx.AddPart(part)
	}
	return ctx
}
