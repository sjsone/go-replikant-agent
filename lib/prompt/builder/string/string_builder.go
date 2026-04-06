package string

import (
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
)

// StringPromptBuilder is a simple prompt builder that concatenates prompts.
type StringPromptBuilder struct {
	raw string
}

// NewStringPromptBuilder creates a new StringPromptBuilder with the given base prompt.
func NewStringPromptBuilder(promptText string) *StringPromptBuilder {
	return &StringPromptBuilder{
		raw: promptText,
	}
}

// Build builds a prompt by concatenating the base prompt with all directive prompts.
func (pb StringPromptBuilder) Build(directives []directive.Directive) prompt.Prompt {
	raw := pb.raw

	for _, d := range directives {
		raw += "\n" + d.GetPrompt().Raw
	}

	return prompt.Prompt{
		Raw: raw,
	}
}
