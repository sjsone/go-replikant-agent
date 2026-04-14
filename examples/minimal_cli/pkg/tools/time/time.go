package time

import (
	"context"
	"fmt"
	"time"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// TimeTool implements the ToolCallable interface for getting current time.
type TimeTool struct {
	tool.Tool
}

// NewTimeTool creates a new time tool.
func NewTimeTool() tool.ToolCallable {
	t := &TimeTool{}

	params, err := tool.SchemaFromStruct[TimeParams]()
	if err != nil {
		panic(err)
	}

	t.Tool = tool.Tool{
		Name:        "get_current_time",
		Description: "Get the current date and time",
		Parameters:  params,
	}
	return t
}

// GetTool returns the tool metadata.
func (t *TimeTool) GetTool() *tool.Tool {
	return &t.Tool
}

// Execute implements ToolCallable.
func (t *TimeTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	now := time.Now()
	return fmt.Sprintf("Current time: %s", now.Format(time.RFC1123)), nil
}
