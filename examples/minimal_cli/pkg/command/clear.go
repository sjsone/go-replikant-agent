package command

import (
	"fmt"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
)

// ClearCommand clears the current conversation
type ClearCommand struct{}

// NewClearCommand creates a new ClearCommand
func NewClearCommand() *ClearCommand {
	return &ClearCommand{}
}

func (c *ClearCommand) GetName() string {
	return "clear"
}

func (c *ClearCommand) GetDescription() string {
	return "Clear the current conversation"
}

func (c *ClearCommand) GetUsage() string {
	return "/clear"
}

func (c *ClearCommand) Execute(args []string, ctx *ExecutionContext) error {
	// Create a new empty context
	newCtx := agentic_context.NewAgentContext()

	// Update the Interactive mode with the new context
	if interactive, ok := ctx.Interactive.(interface {
		SetContext(*agentic_context.AgentContext)
	}); ok {
		interactive.SetContext(newCtx)
	}

	fmt.Fprintln(ctx.Output, "Conversation cleared.")
	return nil
}
