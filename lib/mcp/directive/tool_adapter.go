package directive

import (
	"context"
	"fmt"

	"github.com/sjsone/go-replikant-agent/lib/mcp/client"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// mcpToolAdapter bridges an MCP tool to the framework's ToolCallable interface.
// It implements GetTool() so StaticDirective can auto-extract tool metadata.
type mcpToolAdapter struct {
	toolDef    tool.Tool
	connection client.ServerConnection
}

func newMCPToolAdapter(mcpTool client.MCPTool, connection client.ServerConnection) *mcpToolAdapter {
	params := mcpTool.InputSchema
	if params == nil {
		params = map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}
	}

	return &mcpToolAdapter{
		toolDef: tool.Tool{
			Name:        mcpTool.Name,
			Description: mcpTool.Description,
			Parameters:  params,
		},
		connection: connection,
	}
}

func (a *mcpToolAdapter) GetName() string {
	return a.toolDef.Name
}

func (a *mcpToolAdapter) GetTool() *tool.Tool {
	return &a.toolDef
}

func (a *mcpToolAdapter) Execute(ctx context.Context, args map[string]any) (string, error) {
	result, err := a.connection.CallTool(ctx, a.toolDef.Name, args)
	if err != nil {
		return "", err
	}
	if result.IsError {
		return client.JoinContent(result.Content), fmt.Errorf("MCP tool %s returned an error", a.toolDef.Name)
	}
	return client.JoinContent(result.Content), nil
}
