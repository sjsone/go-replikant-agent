package builder

import (
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
)

// PromptBuilder is an interface for building prompts from directives.
type PromptBuilder interface {
	Build([]directive.Directive) prompt.Prompt
}
