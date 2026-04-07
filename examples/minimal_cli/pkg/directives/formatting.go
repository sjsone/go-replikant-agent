package directives

import (
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// NewTableFormattingDirective creates a Markdown table formatting directive
func NewTableFormattingDirective() directive.Directive {
	return directive.NewStaticDirective(
		"table-format",
		&prompt.Prompt{Raw: "Format the output as a Markdown table. Write the Header of each column in UPPER CASE LETTERS! Do not add text before or after the table. You MUST ONLY print the table. Do not add any markdown code-fences ``` around the table."},
		[]tool.ToolCallable{},
	)
}

// NewCSVFormattingDirective creates a CSV formatting directive
func NewCSVFormattingDirective() directive.Directive {
	return directive.NewStaticDirective(
		"csv-format",
		&prompt.Prompt{Raw: "Format the output as CSV (Comma-Separated Values). Use standard CSV format with commas separating values and newlines separating rows. Do NOT include markdown code blocks."},
		[]tool.ToolCallable{},
	)
}
