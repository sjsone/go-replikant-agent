package directives

import (
	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/tools/time"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// NewTimeDirective creates a time query directive
func NewTimeDirective() directive.Directive {
	return directive.NewStaticDirective(
		"time",
		&prompt.Prompt{Raw: "USE get_current_time when user asks what time it is."},
		[]tool.ToolCallable{time.NewTimeTool()},
	)
}
