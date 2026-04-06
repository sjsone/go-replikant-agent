package agentic_context

import (
	"fmt"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type ContextPartType string

const (
	System ContextPartType = "system"
	User   ContextPartType = "user"
	Agent  ContextPartType = "agent"
	Tool   ContextPartType = "tool"
)

type ContextPart struct {
	Raw         string
	Source      Source
	ToolUse     bool
	Stop        bool
	Cancelled   bool                  // true if this part was created due to cancellation
	ToolCalls   []tool.FunctionCall   `json:"tool_calls,omitempty"`
	ToolResults []tool.FunctionResult `json:"tool_results,omitempty"`

	ConnectedToolCallContextPart       *ContextPart // The ToolCall part this result is from
	ConnectedToolCallSourceContextPart *ContextPart // The part this ToolCall part is from
}

func NewSystemContextPart(raw string) *ContextPart {
	return &ContextPart{
		Raw:    raw,
		Source: SystemSource,
	}
}

func NewAgentContextPart(raw string) *ContextPart {
	return &ContextPart{
		Raw:    raw,
		Source: AgentSource,
	}
}

func NewUserContextPart(raw string) *ContextPart {
	return &ContextPart{
		Raw:    raw,
		Source: UserSource,
	}
}

// NewToolResultContextPart creates a new context part for tool results.
func NewToolResultContextPart(results []tool.FunctionResult) *ContextPart {
	return &ContextPart{
		Raw:         formatToolResults(results),
		Source:      ToolSource,
		ToolResults: results,
	}
}

// IsToolResult returns true if this context part contains tool results.
func (p *ContextPart) IsToolResult() bool {
	return len(p.ToolResults) > 0
}

// formatToolResults formats tool results as a readable string.
func formatToolResults(results []tool.FunctionResult) string {
	s := ""
	for _, r := range results {
		s += fmt.Sprintf("Tool %s: %s\n", r.ID, r.Content)
	}
	return s
}
