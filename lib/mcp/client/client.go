package client

import "context"

// ServerConfig defines how to connect to an MCP server.
// Use Command/Args for stdio transport, or URL for HTTP transport.
type ServerConfig struct {
	Name      string   // Server identifier, used as directive name
	Command   string   // Command to run (stdio transport)
	Args      []string // Arguments for the command (stdio transport)
	Env       []string // Environment variables for the command (stdio transport)
	URL       string   // HTTP endpoint URL (HTTP transport)
	AuthToken string   // Optional Bearer token for HTTP transport authentication
	Prompt    string   // Optional custom system prompt for the directive
}

// MCPTool represents a tool discovered from an MCP server.
type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

// CallResult represents the result of calling an MCP tool.
type CallResult struct {
	Content []string
	IsError bool
}

// ServerConnection abstracts an MCP server connection, translating SDK types
// into framework-friendly types so the directive layer never imports the mcp package.
type ServerConnection interface {
	Connect(ctx context.Context) error
	ListTools(ctx context.Context) ([]MCPTool, error)
	CallTool(ctx context.Context, name string, args map[string]any) (*CallResult, error)
	Close() error
	ServerName() string
}

// NewServerConnection creates a ServerConnection based on the config.
// If URL is set, uses HTTP transport; otherwise uses stdio transport.
func NewServerConnection(config ServerConfig) ServerConnection {
	if config.URL != "" {
		return newHTTPConnection(config)
	}
	return newStdioConnection(config)
}
