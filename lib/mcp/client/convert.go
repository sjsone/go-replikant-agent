package client

import (
	"encoding/json"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// convertTools converts MCP tool definitions to framework-friendly MCPTool types.
func convertTools(tools []*mcp.Tool) []MCPTool {
	result := make([]MCPTool, len(tools))
	for i, t := range tools {
		var schema map[string]any
		if t.InputSchema != nil {
			data, err := json.Marshal(t.InputSchema)
			if err == nil {
				_ = json.Unmarshal(data, &schema)
			}
		}
		result[i] = MCPTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: schema,
		}
	}
	return result
}

// convertResult extracts text content from an MCP CallToolResult.
func convertResult(result *mcp.CallToolResult) *CallResult {
	cr := &CallResult{
		IsError: result.IsError,
	}
	for _, c := range result.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			cr.Content = append(cr.Content, tc.Text)
		}
	}
	return cr
}

// JoinContent joins call result content strings with newlines.
func JoinContent(content []string) string {
	return strings.Join(content, "\n")
}
