package command

import (
	"fmt"
	"sort"
)

// HelpCommand shows help for commands
type HelpCommand struct{}

// NewHelpCommand creates a new HelpCommand
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

func (c *HelpCommand) GetName() string {
	return "help"
}

func (c *HelpCommand) GetDescription() string {
	return "Show help for commands"
}

func (c *HelpCommand) GetUsage() string {
	return "/help [command]"
}

func (c *HelpCommand) Execute(args []string, ctx *ExecutionContext) error {
	// If a specific command is requested, show detailed help
	if len(args) > 0 {
		cmdName := args[0]
		if im, ok := ctx.Interactive.(CommandRegistryGetter); ok {
			if cmd, ok := im.GetCommandRegistry().Get(cmdName); ok {
				fmt.Fprintf(ctx.Output, "Usage: %s\n", cmd.GetUsage())
				fmt.Fprintf(ctx.Output, "\n%s\n", cmd.GetDescription())
				return nil
			}
			return fmt.Errorf("unknown command: %s", cmdName)
		}
	}

	// Show general help
	if im, ok := ctx.Interactive.(CommandRegistryGetter); ok {
		cmds := im.GetCommandRegistry().List()

		// Sort commands by name
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].GetName() < cmds[j].GetName()
		})

		fmt.Fprintln(ctx.Output, "Available commands:")
		fmt.Fprintln(ctx.Output)
		for _, cmd := range cmds {
			fmt.Fprintf(ctx.Output, "  %-12s %s\n", "/"+cmd.GetName(), cmd.GetDescription())
		}
		fmt.Fprintln(ctx.Output)
		fmt.Fprintln(ctx.Output, "Use /help <command> for more information on a specific command.")
		return nil
	}

	return nil
}
