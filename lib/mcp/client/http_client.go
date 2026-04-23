package client

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type httpConnection struct {
	config  ServerConfig
	client  *mcp.Client
	session *mcp.ClientSession
}

func newHTTPConnection(config ServerConfig) *httpConnection {
	return &httpConnection{config: config}
}

// authTransport wraps an http.RoundTripper to inject an Authorization header.
type authTransport struct {
	base   http.RoundTripper
	token  string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}
	return t.base.RoundTrip(req)
}

func (c *httpConnection) Connect(ctx context.Context) error {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "replikant-agent",
		Version: "0.1.0",
	}, nil)

	transport := &mcp.StreamableClientTransport{Endpoint: c.config.URL}
	if c.config.AuthToken != "" {
		transport.HTTPClient = &http.Client{
			Transport: &authTransport{
				base:  http.DefaultTransport,
				token: c.config.AuthToken,
			},
		}
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return err
	}

	c.client = client
	c.session = session
	return nil
}

func (c *httpConnection) ListTools(ctx context.Context) ([]MCPTool, error) {
	result, err := c.session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, err
	}
	return convertTools(result.Tools), nil
}

func (c *httpConnection) CallTool(ctx context.Context, name string, args map[string]any) (*CallResult, error) {
	result, err := c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}
	return convertResult(result), nil
}

func (c *httpConnection) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}

func (c *httpConnection) ServerName() string {
	return c.config.Name
}
