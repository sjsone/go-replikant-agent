package mode

import (
	"context"
	"fmt"
	"os"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/cli"
	"github.com/sjsone/go-replikant-agent/lib/router"
	lib_session "github.com/sjsone/go-replikant-agent/lib/session"
)

func RunSingleTurnMode(
	agenticSession *lib_session.AgenticSession,
	initialContext *agentic_context.AgentContext,
	prompt string,
	router router.Router,
) {
	ctx := context.Background()

	// Create a simple delegate for output
	delegate := cli.NewPrettyPrintDelegate()
	agenticSession.SetDelegate(delegate)
	router.SetDelegate(delegate)

	// Create user context part
	userPart := agentic_context.NewUserContextPart(prompt)

	// Process the prompt
	if err := agenticSession.ProcessContextPart(ctx, userPart); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Ensure final newline after streaming
	fmt.Println()
}
