package command

import "io"

// Command represents a user command (e.g., /exit, /save)
type Command interface {
	GetName() string        // Command name without '/'
	GetDescription() string // Short description
	GetUsage() string       // Usage info
	Execute(args []string, ctx *ExecutionContext) error
}

// CommandRegistryGetter is implemented by types that provide access to a CommandRegistry
type CommandRegistryGetter interface {
	GetCommandRegistry() *CommandRegistry
}

// ExecutionContext provides context for command execution
type ExecutionContext struct {
	Interactive any // *Interactive from pkg/cli (avoiding circular import)
	Output      io.Writer
}

// CommandRegistry manages available commands
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry creates a new CommandRegistry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

// Register registers a command in the registry
func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.GetName()] = cmd
}

// Get retrieves a command by name
func (r *CommandRegistry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// List returns all registered commands
func (r *CommandRegistry) List() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}
