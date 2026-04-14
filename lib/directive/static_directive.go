package directive

import (
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type StaticDirective struct {
	name      string
	prompt    *prompt.Prompt
	callables []tool.ToolCallable
}

func NewStaticDirective(name string, prompt *prompt.Prompt, callables []tool.ToolCallable) *StaticDirective {
	return &StaticDirective{
		name:      name,
		prompt:    prompt,
		callables: callables,
	}
}

func (d *StaticDirective) GetName() string {
	return d.name
}

func (d *StaticDirective) GetPrompt() *prompt.Prompt {
	return d.prompt
}

func (d *StaticDirective) GetToolCallables() []tool.ToolCallable {
	return d.callables
}
