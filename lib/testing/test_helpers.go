package testing

import (
	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// NewTestAgentContext creates a new AgentContext for testing.
func NewTestAgentContext() *agentic_context.AgentContext {
	return agentic_context.NewAgentContext()
}

// NewTestContextPart creates a ContextPart with the given source and content.
func NewTestContextPart(source agentic_context.Source, content string) *agentic_context.ContextPart {
	return &agentic_context.ContextPart{
		Raw:    content,
		Source: source,
	}
}

// NewTestTool creates a new Tool for testing.
func NewTestTool(name, description string) *tool.Tool {
	return &tool.Tool{
		Name:        name,
		Description: description,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Test input parameter",
				},
			},
		},
	}
}

// NewTestDirective creates a StaticDirective for testing.
func NewTestDirective(name, promptText string) *directive.StaticDirective {
	p := &prompt.Prompt{
		Raw: promptText,
	}
	return directive.NewStaticDirective(name, p, []tool.ToolCallable{})
}

// NewTestUserContextPart creates a user context part for testing.
func NewTestUserContextPart(content string) *agentic_context.ContextPart {
	return agentic_context.NewUserContextPart(content)
}

// NewTestAgentContextPart creates an agent context part for testing.
func NewTestAgentContextPart(content string) *agentic_context.ContextPart {
	return agentic_context.NewAgentContextPart(content)
}

// NewTestSystemContextPart creates a system context part for testing.
func NewTestSystemContextPart(content string) *agentic_context.ContextPart {
	return agentic_context.NewSystemContextPart(content)
}

// NewTestToolResultContextPart creates a tool result context part for testing.
func NewTestToolResultContextPart(results []tool.FunctionResult) *agentic_context.ContextPart {
	return agentic_context.NewToolResultContextPart(results)
}
