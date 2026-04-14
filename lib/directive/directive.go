package directive

import (
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type Directive interface {
	GetName() string
	GetPrompt() *prompt.Prompt
	GetToolCallables() []tool.ToolCallable
}
