package directive

import (
	"context"
	"fmt"
	"strings"

	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/mcp/client"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// MCPDirective wraps a StaticDirective with tools discovered from an MCP server.
// One MCPDirective per MCP server. Use Close() to clean up the connection.
type MCPDirective struct {
	inner      *directive.StaticDirective
	connection client.ServerConnection
}

// NewMCPDirective connects to an MCP server, discovers its tools, and builds a directive.
func NewMCPDirective(ctx context.Context, config client.ServerConfig) (*MCPDirective, error) {
	connection := client.NewServerConnection(config)

	if err := connection.Connect(ctx); err != nil {
		return nil, fmt.Errorf("MCP connect to %s: %w", config.Name, err)
	}

	mcpTools, err := connection.ListTools(ctx)
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("MCP list tools from %s: %w", config.Name, err)
	}

	callables := make([]tool.ToolCallable, len(mcpTools))
	for i, mt := range mcpTools {
		callables[i] = newMCPToolAdapter(mt, connection)
	}

	dirPrompt := buildPrompt(config, mcpTools)
	inner := directive.NewStaticDirective(config.Name, dirPrompt, callables)

	return &MCPDirective{
		inner:      inner,
		connection: connection,
	}, nil
}

func (d *MCPDirective) Directive() directive.Directive {
	return d.inner
}

func (d *MCPDirective) Close() error {
	return d.connection.Close()
}

func buildPrompt(config client.ServerConfig, tools []client.MCPTool) *prompt.Prompt {
	if config.Prompt != "" {
		return &prompt.Prompt{Raw: config.Prompt}
	}

	var sb strings.Builder
	sb.WriteString("You have access to the following tools from the ")
	sb.WriteString(config.Name)
	sb.WriteString(" MCP server:\n\n")
	for _, t := range tools {
		sb.WriteString("- ")
		sb.WriteString(t.Name)
		if t.Description != "" {
			sb.WriteString(": ")
			sb.WriteString(t.Description)
		}
		sb.WriteString("\n")
	}

	return &prompt.Prompt{Raw: sb.String()}
}
