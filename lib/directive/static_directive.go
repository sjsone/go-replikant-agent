package directive

import (
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type StaticDirective struct {
	name      string
	prompt    *prompt.Prompt
	tools     []*tool.Tool
	callables []tool.ToolCallable
}

func NewStaticDirective(name string, prompt *prompt.Prompt, callables []tool.ToolCallable) *StaticDirective {
	// Extract tool metadata from callables
	tools := make([]*tool.Tool, 0, len(callables))
	for _, c := range callables {
		// Try to get tool metadata from the callable
		// Callables that embed tool.Tool will have the fields accessible
		if t, ok := c.(interface{ GetTool() *tool.Tool }); ok {
			tools = append(tools, t.GetTool())
		}
	}

	return &StaticDirective{
		name:      name,
		prompt:    prompt,
		tools:     tools,
		callables: callables,
	}
}

func (d *StaticDirective) GetName() string {
	return d.name
}

func (d *StaticDirective) GetPrompt() *prompt.Prompt {
	return d.prompt
}

func (d *StaticDirective) GetTools() []*tool.Tool {
	return d.tools
}

func (d *StaticDirective) GetToolCallables() []tool.ToolCallable {
	return d.callables
}
