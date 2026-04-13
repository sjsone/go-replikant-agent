package wait

import (
	"context"
	"fmt"
	"time"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

type WaitTool struct {
	tool.Tool
}

type WaitParams struct {
	Seconds int `json:"seconds" jsonschema:"description=Number of seconds to wait (1-20),minimum=1,maximum=20"`
}

func NewWaitTool() tool.ToolCallable {
	t := &WaitTool{}

	params, err := tool.SchemaFromStruct[WaitParams]()
	if err != nil {
		panic(err)
	}

	t.Tool = tool.Tool{
		Name:        "wait",
		Description: "Wait for a specified number of seconds (1-20). Use this when you receive rate-limit errors (e.g. TOO MANY REQUESTS, 429) and need to back off before retrying.",
		Parameters:  params,
	}
	return t
}

func (t *WaitTool) GetTool() *tool.Tool {
	return &t.Tool
}

func (t *WaitTool) GetName() string {
	return t.Tool.Name
}

func (t *WaitTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	seconds := 1
	if v, ok := args["seconds"].(float64); ok {
		seconds = int(v)
	}
	if seconds < 1 {
		seconds = 1
	}
	if seconds > 20 {
		seconds = 20
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Duration(seconds) * time.Second):
		return fmt.Sprintf("Waited %d seconds.", seconds), nil
	}
}
