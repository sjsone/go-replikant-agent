package mode

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	command "github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/command"
	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/cli"
	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/router"
	libsession "github.com/sjsone/go-replikant-agent/lib/session"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// InteractiveMode implements SessionDelegate for REPL behavior
type InteractiveMode struct {
	session         *libsession.AgenticSession
	context         *agentic_context.AgentContext
	commandRegistry *command.CommandRegistry
	allDirectives   *[]directive.Directive
	output          io.Writer

	// Cancellation
	cancelFunc        context.CancelFunc
	cancelMutex       sync.Mutex
	requestInProgress bool
	requestMutex      sync.Mutex

	// Compose PrettyPrintDelegate for output formatting
	outputDelegate *cli.PrettyPrintDelegate

	// State
	shouldExit bool
	prompt     string
}

// NewInteractiveMode creates a new InteractiveMode
func NewInteractiveMode(
	sess *libsession.AgenticSession,
	initialCtx *agentic_context.AgentContext,
	allDirectives *[]directive.Directive,
) *InteractiveMode {
	im := &InteractiveMode{
		session:         sess,
		context:         initialCtx,
		commandRegistry: command.NewCommandRegistry(),
		allDirectives:   allDirectives,
		output:          os.Stdout,
		outputDelegate:  cli.NewPrettyPrintDelegate(),
		prompt:          "> ",
	}

	// Register commands
	im.commandRegistry.Register(command.NewExitCommand())
	im.commandRegistry.Register(command.NewClearCommand())
	im.commandRegistry.Register(command.NewHelpCommand())

	return im
}

// Run starts the REPL loop
func (im *InteractiveMode) Run(ctx context.Context) error {
	// Show welcome message
	im.showWelcome()

	// Handle signals in a dedicated goroutine so Ctrl-C works
	// even while a request is in progress (which blocks the main select).
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	exitChan := make(chan struct{})
	go func() {
		for range sigChan {
			im.requestMutex.Lock()
			inProgress := im.requestInProgress
			im.requestMutex.Unlock()

			if inProgress {
				im.CancelCurrentRequest()
				fmt.Fprintln(im.output, "\n[Cancelled]")
			} else {
				close(exitChan)
				return
			}
		}
	}()

	// Input loop
	scanner := bufio.NewScanner(os.Stdin)
	for !im.shouldExit {
		// Show prompt
		fmt.Fprint(im.output, im.prompt)

		inputChan := make(chan string)
		errChan := make(chan error, 1)

		go func() {
			if scanner.Scan() {
				inputChan <- scanner.Text()
			} else {
				if err := scanner.Err(); err != nil {
					errChan <- err
				} else {
					// EOF
					close(inputChan)
				}
			}
		}()

		select {
		case line, ok := <-inputChan:
			if !ok {
				// EOF, exit normally
				return nil
			}
			line = strings.TrimSpace(line)
			if line != "" {
				im.processInput(ctx, line)
			}
		case err := <-errChan:
			return err
		case <-exitChan:
			im.Exit()
			fmt.Fprintln(im.output, "Goodbye!")
			return nil
		}
	}

	return nil
}

// showWelcome displays a welcome message
func (im *InteractiveMode) showWelcome() {
	reset := "\033[0m"
	bold := "\033[1m"
	yellow := "\033[33m"
	fmt.Fprintln(im.output, bold+yellow+" Replikant Agent - Minimal CLI Example"+reset)
	fmt.Fprintln(im.output, "")

	fmt.Fprintln(im.output, "Type your message or /help for commands.")
	fmt.Fprintln(im.output)
}

// processInput processes a line of input
func (im *InteractiveMode) processInput(ctx context.Context, line string) {
	// Check if it's a command
	if strings.HasPrefix(line, "/") {
		im.executeCommand(line)
		return
	}

	// Otherwise, treat as chat message
	im.processChatMessage(ctx, line)
}

// processChatMessage processes a chat message
func (im *InteractiveMode) processChatMessage(ctx context.Context, message string) {
	// Mark request as in progress
	im.setRequestInProgress(true)
	defer im.setRequestInProgress(false)

	// Create cancellable context for this request
	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	im.setCancelFunc(cancel)

	// Create user context part
	userPart := agentic_context.NewUserContextPart(message)

	// Process the part
	err := im.session.ProcessContextPart(reqCtx, userPart)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			fmt.Fprintf(im.output, "\nError: %v\n", err)
		}
	}
}

// executeCommand executes a command
func (im *InteractiveMode) executeCommand(line string) {
	// Parse command and arguments
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	cmdName := strings.TrimPrefix(parts[0], "/")
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	// Look up command
	cmd, ok := im.commandRegistry.Get(cmdName)
	if !ok {
		fmt.Fprintf(im.output, "Unknown command: %s\n", cmdName)
		fmt.Fprintln(im.output, "Use /help for available commands.")
		return
	}

	// Create execution context
	execCtx := &command.ExecutionContext{
		Interactive: im,
		Output:      im.output,
	}

	// Execute command
	if err := cmd.Execute(args, execCtx); err != nil {
		fmt.Fprintf(im.output, "Error executing command: %v\n", err)
	}
}

// CancelCurrentRequest cancels the current request
func (im *InteractiveMode) CancelCurrentRequest() {
	im.cancelMutex.Lock()
	defer im.cancelMutex.Unlock()

	if im.cancelFunc != nil {
		im.cancelFunc()
	}
}

// Exit sets the shouldExit flag to true
func (im *InteractiveMode) Exit() {
	im.shouldExit = true
}

// setCancelFunc sets the current cancel function
func (im *InteractiveMode) setCancelFunc(cancel context.CancelFunc) {
	im.cancelMutex.Lock()
	defer im.cancelMutex.Unlock()
	im.cancelFunc = cancel
}

// setRequestInProgress sets the request in progress state
func (im *InteractiveMode) setRequestInProgress(inProgress bool) {
	im.requestMutex.Lock()
	defer im.requestMutex.Unlock()
	im.requestInProgress = inProgress
}

// GetContext returns the current context
func (im *InteractiveMode) GetContext() *agentic_context.AgentContext {
	return im.context
}

// SetContext sets the current context
func (im *InteractiveMode) SetContext(ctx *agentic_context.AgentContext) {
	im.context = ctx
}

// GetCommandRegistry returns the command registry
func (im *InteractiveMode) GetCommandRegistry() *command.CommandRegistry {
	return im.commandRegistry
}

// MARK: Delegate implementations

// MARK: Session Delegate
func (im *InteractiveMode) SessionOnPartAdded(part *agentic_context.ContextPart) {
	im.outputDelegate.SessionOnPartAdded(part)
}

func (im *InteractiveMode) SessionOnToolCallsReceived(calls []tool.FunctionCall) {
	im.outputDelegate.SessionOnToolCallsReceived(calls)
}

func (im *InteractiveMode) SessionOnToolExecuted(call tool.FunctionCall, result tool.FunctionResult) {
	im.outputDelegate.SessionOnToolExecuted(call, result)
}

func (im *InteractiveMode) SessionOnStreamingChunk(chunk string) {
	im.outputDelegate.SessionOnStreamingChunk(chunk)
}

func (im *InteractiveMode) SessionOnRequestSent(messages []connector.Message, directives []directive.Directive) {
	im.outputDelegate.SessionOnRequestSent(messages, directives)
}

func (im *InteractiveMode) SessionOnLoopIteration(iteration int) {
	im.outputDelegate.SessionOnLoopIteration(iteration)
}

func (im *InteractiveMode) SessionOnLoopEnd() {
	im.outputDelegate.SessionOnLoopEnd()
}

// MARK: Router Delegate
func (im *InteractiveMode) RouterOnRoutingDecision(decision router.RoutingDecision, allOptions []*router.RoutingOption, activeOptions []*router.RoutingOption) {
	im.outputDelegate.RouterOnRoutingDecision(decision, allOptions, activeOptions)
}
