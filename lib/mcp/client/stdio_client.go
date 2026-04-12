package client

import (
	"context"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type stdioConnection struct {
	config  ServerConfig
	client  *mcp.Client
	session *mcp.ClientSession
}

func newStdioConnection(config ServerConfig) *stdioConnection {
	return &stdioConnection{config: config}
}

func (c *stdioConnection) Connect(ctx context.Context) error {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "replikant-agent",
		Version: "0.1.0",
	}, nil)

	cmd := exec.Command(c.config.Command, c.config.Args...)
	if len(c.config.Env) > 0 {
		cmd.Env = c.config.Env
	}

	transport := &mcp.CommandTransport{Command: cmd}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return err
	}

	c.client = client
	c.session = session
	return nil
}

func (c *stdioConnection) ListTools(ctx context.Context) ([]MCPTool, error) {
	result, err := c.session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, err
	}
	return convertTools(result.Tools), nil
}

func (c *stdioConnection) CallTool(ctx context.Context, name string, args map[string]any) (*CallResult, error) {
	result, err := c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}
	return convertResult(result), nil
}

func (c *stdioConnection) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}

func (c *stdioConnection) ServerName() string {
	return c.config.Name
}
