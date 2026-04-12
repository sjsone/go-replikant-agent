package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/router"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// ANSI color codes for terminal output
const (
	reset = "\033[0m"
	bold  = "\033[1m"
	dim   = "\033[2m"

	// Colors
	gray    = "\033[90m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	// Bright colors
	brightGreen  = "\033[92m"
	brightYellow = "\033[93m"
	brightBlue   = "\033[94m"
)

// PrettyPrintDelegate provides a pretty CLI output behavior with colors and formatting.
type PrettyPrintDelegate struct {
	firstChunk     bool
	inToolBlock    bool
	iteration      int
	showIterations bool
	pendingNewline bool
}

// NewPrettyPrintDelegate creates a new PrettyPrintDelegate.
func NewPrettyPrintDelegate() *PrettyPrintDelegate {
	return &PrettyPrintDelegate{
		firstChunk:  true,
		inToolBlock: false,
		iteration:   0,
	}
}

// ShowIterations enables display of loop iteration numbers.
func (d *PrettyPrintDelegate) ShowIterations(show bool) {
	d.showIterations = show
}

func (d *PrettyPrintDelegate) SessionOnPartAdded(part *agentic_context.ContextPart) {
	// Add newline after streaming if needed
	if d.pendingNewline {
		fmt.Println()
		d.pendingNewline = false
	}

	// Handle cancelled parts
	if part.Cancelled {
		fmt.Printf("%s%s⚠️  Response cancelled%s\n", yellow, bold, reset)
		if part.Raw != "" {
			fmt.Printf("%sPartial response: %s%s\n", dim, part.Raw, reset)
		}
		return
	}

	if part.Source.IsSystem() {
		return
	}
}

func (d *PrettyPrintDelegate) SessionOnToolCallsReceived(calls []tool.FunctionCall) {
	if len(calls) == 0 {
		return
	}

	// Add newline after streaming if needed
	if d.pendingNewline {
		fmt.Println()
		d.pendingNewline = false
	}

	fmt.Printf("\n%s", dim)
	fmt.Println("┌──────────────────────────────────────────────────────────")
	for _, call := range calls {
		// Format args for display
		argsJson, _ := json.Marshal(call.Arguments)
		argsStr := string(argsJson)

		// Truncate long args
		if len(argsStr) > 60 {
			argsStr = argsStr[:60] + "..."
		}

		fmt.Printf("%s├", dim)
		fmt.Printf(" %sTool call:%s %s%s(%s)%s\n", brightBlue, reset, cyan, call.Name, argsStr, reset)
	}
	fmt.Printf("%s└─────────────────────────────────────────────────────────%s", dim, reset)

	d.inToolBlock = true
}

func (d PrettyPrintDelegate) getColoredCheckIcon(ok bool) (string, string) {
	icon := "✓"
	color := green
	if !ok {
		icon = "✗"
		color = red
	}

	return icon, color
}

func (d *PrettyPrintDelegate) SessionOnToolExecuted(call tool.FunctionCall, result tool.FunctionResult) {
	icon, color := d.getColoredCheckIcon(!result.Error)

	// Format result content
	content := strings.TrimSpace(result.Content)
	if content == "" {
		content = "(no output)"
	}

	// Truncate long results
	maxLen := 200
	if len(content) > maxLen {
		content = content[:maxLen] + "..."
	}

	// Replace newlines with spaces for single-line display
	content = strings.ReplaceAll(content, "\n", " ")

	fmt.Printf("\n%s  %s%s %s%s%s%s", dim, color, icon, reset, dim, content, reset)
}

func (d *PrettyPrintDelegate) SessionOnStreamingChunk(chunk string) {
	if d.firstChunk {
		// End tool block with newline if we were in one
		if d.inToolBlock {
			fmt.Println()
			d.inToolBlock = false
		}

		// Show iteration indicator if enabled
		if d.showIterations && d.iteration > 0 {
			fmt.Printf("\n%s[%d]%s ", dim, d.iteration, reset)
		}
		fmt.Printf("\n%s%sAgent:%s ", bold, brightGreen, reset)
		d.firstChunk = false
	}
	fmt.Print(chunk)
	d.pendingNewline = true
}

func (d *PrettyPrintDelegate) SessionOnRequestSent(messages []connector.Message, directives []directive.Directive) {
	// Reset first chunk flag for new request
	if !d.firstChunk {
		d.firstChunk = true
	}
}

func (d *PrettyPrintDelegate) SessionOnLoopIteration(iteration int) {
	d.iteration = iteration
}

func (d *PrettyPrintDelegate) SessionOnLoopEnd() {
	if d.pendingNewline {
		fmt.Println()
		d.pendingNewline = false
	}
	if d.inToolBlock {
		fmt.Println()
		d.inToolBlock = false
	}
}

func (d *PrettyPrintDelegate) RouterPreparedRouting(allOptions []*router.RoutingOption) {
	fmt.Printf("%s╔Routing Options══════════%s\n", dim, reset)

	for _, o := range allOptions {
		fmt.Printf(" -- %s \n", o.Name)
		if len(o.ToolNames) > 0 {
			fmt.Printf("    Tools: \n")
			for _, n := range o.ToolNames {
				fmt.Printf("      - %s \n", n)
			}
		}
	}

	fmt.Printf("\n")
}

// RouterOnRoutingDecision is called when RouterMultiplexer makes a routing decision.
func (d *PrettyPrintDelegate) RouterOnRoutingDecision(decision router.RoutingDecision, allOptions []*router.RoutingOption, activeOptions []*router.RoutingOption) {
	// Add newline after streaming if needed
	if d.pendingNewline {
		fmt.Println()
		d.pendingNewline = false
	}

	fmt.Printf("%s╔Routing══════════%s\n", dim, reset)
	fmt.Printf("%s╟─%s[options]%s%s: %v %s\n", dim, brightBlue, reset, dim, decision.SelectedIDs, reset)
	if decision.Reasoning != "" {
		// Truncate long reasoning
		reasoning := decision.Reasoning
		reasoningTruncateLength := 1500
		if len(reasoning) > reasoningTruncateLength {
			reasoning = reasoning[:reasoningTruncateLength] + "..."
		}
		words := strings.Fields(reasoning)
		line := ""
		for _, word := range words {
			candidate := line + " " + word
			if len(candidate) > 70 && line != "" {
				fmt.Printf("%s║ %s%s%s\n", dim, dim, line, reset)
				line = word
			} else {
				line = candidate
			}
		}
		if line != "" {
			fmt.Printf("%s║ %s%s%s\n", dim, dim, strings.TrimSpace(line), reset)
		}
	}
	fmt.Printf("%s╚═════════════════%s\n", dim, reset)

	// Display active directives with tools
	if len(activeOptions) > 0 {
		// Calculate box width
		header := " Directives "
		boxWidth := len(header) + 2
		for _, opt := range activeOptions {
			lineLen := len(opt.Name) + 4 // "│ []: "
			for _, tn := range opt.ToolNames {
				lineLen += len(tn) + 4 // ""tool_name", "
				if lineLen > boxWidth {
					boxWidth = lineLen
				}
			}
		}
		// Ensure minimum width
		if boxWidth < 13 {
			boxWidth = 13
		}

		fmt.Println()
		// Print top border
		fmt.Printf("%s╭Directives%s%s\n", dim, strings.Repeat("─", boxWidth-2), reset)

		// Print each directive with its tools
		for opt_index, opt := range activeOptions {
			//
			//
			c := "├"
			if opt_index == 0 {
				c = "╭"
			}

			fmt.Printf("%s│%s%s[%s]%s%s:\n", dim, c, cyan, opt.Name, reset, dim)

			optPrompt := strings.Trim(opt.Text, "")
			optPrompt = strings.ReplaceAll(optPrompt, "\n", " ")
			optPromptTruncateLength := 40
			if len(optPrompt) > optPromptTruncateLength {
				optPrompt = optPrompt[:optPromptTruncateLength] + "..."
			}

			if optPrompt != "" {
				fmt.Printf("%s││├─Prompt: '%s'%s\n", dim, optPrompt, reset)
			} else {
				promptIcon, promptIconColor := d.getColoredCheckIcon(false)
				fmt.Printf("%s││├─Prompt: %s%s%s\n", dim, promptIconColor, promptIcon, reset)
			}

			if len(opt.ToolNames) > 0 {
				fmt.Printf("%s││╰─Tools: %v%s\n", dim, opt.ToolNames, reset)
			} else {
				toolsIcon, toolsIconColor := d.getColoredCheckIcon(false)
				fmt.Printf("%s││╰─Tools: %s%s%s\n", dim, toolsIconColor, toolsIcon, reset)
			}
		}

		// Print bottom border
		fmt.Printf("%s╰╯%s\n", dim, reset)
	}
}
