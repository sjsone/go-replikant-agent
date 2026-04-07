package command

import (
	"fmt"
)

// ExitCommand exits the interactive session
type ExitCommand struct{}

// NewExitCommand creates a new ExitCommand
func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

func (c *ExitCommand) GetName() string {
	return "exit"
}

func (c *ExitCommand) GetDescription() string {
	return "Exit the interactive session"
}

func (c *ExitCommand) GetUsage() string {
	return "/exit"
}

func (c *ExitCommand) Execute(args []string, ctx *ExecutionContext) error {
	// Get the Interactive mode from context and call Exit
	if interactive, ok := ctx.Interactive.(interface{ Exit() }); ok {
		interactive.Exit()
	}
	fmt.Fprintln(ctx.Output, "Goodbye!")
	return nil
}
