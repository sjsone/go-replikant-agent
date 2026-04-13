package directives

import (
	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/tools/wait"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// NewWaitDirective creates a directive with a wait/backoff tool, always available for rate-limit recovery.
func NewWaitDirective() directive.Directive {
	return directive.NewStaticDirective(
		"wait",
		&prompt.Prompt{Raw: `
RATE-LIMIT AND BACKOFF HANDLING:

When you receive rate-limit errors such as:
- "TOO MANY REQUESTS"
- HTTP 429
- "rate limit exceeded"
- "too many requests"

Use the wait tool to back off before retrying. Recommended wait times:
- First retry: 2-5 seconds
- Subsequent retries: 5-10 seconds
- Extended backoff: 10-20 seconds

Always explain to the user that you are waiting due to rate limiting.
`},
		[]tool.ToolCallable{wait.NewWaitTool()},
	)
}
