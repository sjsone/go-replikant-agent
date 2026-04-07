package mode

import (
	"context"
	"fmt"
	"os"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/router"
	lib_session "github.com/sjsone/go-replikant-agent/lib/session"
)

func RunInteractiveMode(
	agenticSession *lib_session.AgenticSession,
	initialContext *agentic_context.AgentContext,
	router router.Router,
	allDirectives *[]directive.Directive,
) {
	// Create interactive mode
	interactiveMode := NewInteractiveMode(agenticSession, initialContext, allDirectives)

	// Set the interactive mode as the session delegate
	agenticSession.SetDelegate(interactiveMode)

	router.SetDelegate(interactiveMode)

	// Create context for the REPL
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the REPL
	if err := interactiveMode.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
