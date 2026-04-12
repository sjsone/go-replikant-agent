package directive

import (
	"context"
	"errors"
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/mcp/client"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// mockConnection implements client.ServerConnection for testing.
type mockConnection struct {
	name  string
	tools []client.MCPTool
	// callTool tracks the last call and returns a configured result
	callToolFn func(name string, args map[string]any) (*client.CallResult, error)
	closed     bool
}

func (m *mockConnection) Connect(ctx context.Context) error { return nil }
func (m *mockConnection) ListTools(ctx context.Context) ([]client.MCPTool, error) {
	return m.tools, nil
}
func (m *mockConnection) CallTool(ctx context.Context, name string, args map[string]any) (*client.CallResult, error) {
	if m.callToolFn != nil {
		return m.callToolFn(name, args)
	}
	return &client.CallResult{Content: []string{"ok"}}, nil
}
func (m *mockConnection) Close() error { m.closed = true; return nil }
func (m *mockConnection) ServerName() string { return m.name }

// --- mcpToolAdapter tests ---

func TestMCPToolAdapter_GetName(t *testing.T) {
	mt := client.MCPTool{Name: "test_tool", Description: "A test tool"}
	conn := &mockConnection{name: "test-server"}
	adapter := newMCPToolAdapter(mt, conn)

	if got := adapter.GetName(); got != "test_tool" {
		t.Errorf("GetName() = %q, want %q", got, "test_tool")
	}
}

func TestMCPToolAdapter_GetTool(t *testing.T) {
	mt := client.MCPTool{
		Name:        "my_tool",
		Description: "Does things",
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{"x": map[string]any{"type": "string"}}},
	}
	conn := &mockConnection{name: "test"}
	adapter := newMCPToolAdapter(mt, conn)

	got := adapter.GetTool()
	if got.Name != "my_tool" {
		t.Errorf("GetTool().Name = %q, want %q", got.Name, "my_tool")
	}
	if got.Description != "Does things" {
		t.Errorf("GetTool().Description = %q, want %q", got.Description, "Does things")
	}
	if got.Parameters == nil {
		t.Error("GetTool().Parameters is nil")
	}
}

func TestMCPToolAdapter_GetTool_NilSchema(t *testing.T) {
	mt := client.MCPTool{Name: "no_schema", Description: "No input schema"}
	conn := &mockConnection{name: "test"}
	adapter := newMCPToolAdapter(mt, conn)

	got := adapter.GetTool()
	if got.Parameters == nil {
		t.Error("GetTool().Parameters should have default schema when InputSchema is nil")
	}
	if got.Parameters["type"] != "object" {
		t.Errorf("default schema type = %v, want %q", got.Parameters["type"], "object")
	}
}

func TestMCPToolAdapter_ImplementsToolCallable(t *testing.T) {
	// Verify the adapter satisfies the interfaces we need
	mt := client.MCPTool{Name: "t"}
	conn := &mockConnection{name: "s"}
	adapter := newMCPToolAdapter(mt, conn)

	var _ tool.ToolCallable = adapter
	var _ interface{ GetTool() *tool.Tool } = adapter
}

func TestMCPToolAdapter_Execute(t *testing.T) {
	mt := client.MCPTool{Name: "echo", Description: "Echoes input"}
	conn := &mockConnection{
		name: "test",
		callToolFn: func(name string, args map[string]any) (*client.CallResult, error) {
			if name != "echo" {
				t.Errorf("CallTool name = %q, want %q", name, "echo")
			}
			return &client.CallResult{Content: []string{"hello world"}}, nil
		},
	}
	adapter := newMCPToolAdapter(mt, conn)

	result, err := adapter.Execute(context.Background(), map[string]any{"msg": "hi"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result != "hello world" {
		t.Errorf("Execute() = %q, want %q", result, "hello world")
	}
}

func TestMCPToolAdapter_Execute_MultiContent(t *testing.T) {
	mt := client.MCPTool{Name: "multi"}
	conn := &mockConnection{
		name: "test",
		callToolFn: func(name string, args map[string]any) (*client.CallResult, error) {
			return &client.CallResult{Content: []string{"line1", "line2", "line3"}}, nil
		},
	}
	adapter := newMCPToolAdapter(mt, conn)

	result, err := adapter.Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	want := "line1\nline2\nline3"
	if result != want {
		t.Errorf("Execute() = %q, want %q", result, want)
	}
}

func TestMCPToolAdapter_Execute_ToolError(t *testing.T) {
	mt := client.MCPTool{Name: "fail_tool"}
	conn := &mockConnection{
		name: "test",
		callToolFn: func(name string, args map[string]any) (*client.CallResult, error) {
			return &client.CallResult{Content: []string{"something went wrong"}, IsError: true}, nil
		},
	}
	adapter := newMCPToolAdapter(mt, conn)

	_, err := adapter.Execute(context.Background(), nil)
	if err == nil {
		t.Fatal("Execute() should return error when tool returns IsError")
	}
}

func TestMCPToolAdapter_Execute_ConnectionError(t *testing.T) {
	mt := client.MCPTool{Name: "broken"}
	conn := &mockConnection{
		name: "test",
		callToolFn: func(name string, args map[string]any) (*client.CallResult, error) {
			return nil, errors.New("connection refused")
		},
	}
	adapter := newMCPToolAdapter(mt, conn)

	_, err := adapter.Execute(context.Background(), nil)
	if err == nil {
		t.Fatal("Execute() should return error when connection fails")
	}
}

// --- MCPDirective tests ---

func TestNewMCPDirective(t *testing.T) {
	tools := []client.MCPTool{
		{Name: "tool_a", Description: "Tool A"},
		{Name: "tool_b", Description: "Tool B", InputSchema: map[string]any{"type": "object"}},
	}
	conn := &mockConnection{name: "my-server", tools: tools}

	// We can't use NewMCPDirective directly because it creates a connection internally.
	// Instead, test the inner parts by constructing manually.
	callables := make([]tool.ToolCallable, len(tools))
	for i, mt := range tools {
		callables[i] = newMCPToolAdapter(mt, conn)
	}

	dir := buildTestDirective(conn, tools)

	impl := dir.Directive()
	if impl.GetName() != "my-server" {
		t.Errorf("GetName() = %q, want %q", impl.GetName(), "my-server")
	}
	if len(impl.GetTools()) != 2 {
		t.Errorf("GetTools() count = %d, want 2", len(impl.GetTools()))
	}
	if len(impl.GetToolCallables()) != 2 {
		t.Errorf("GetToolCallables() count = %d, want 2", len(impl.GetToolCallables()))
	}

	// Verify tool names
	names := make(map[string]bool)
	for _, tc := range impl.GetToolCallables() {
		names[tc.GetName()] = true
	}
	if !names["tool_a"] || !names["tool_b"] {
		t.Errorf("missing tool names, got %v", names)
	}

	// Verify prompt mentions the server
	prompt := impl.GetPrompt()
	if prompt == nil {
		t.Fatal("GetPrompt() is nil")
	}

	// Close should work
	if err := dir.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
	if !conn.closed {
		t.Error("connection was not closed")
	}

	// Suppress unused warning
	_ = callables
}

func TestBuildPrompt_Default(t *testing.T) {
	config := client.ServerConfig{Name: "test-server"}
	tools := []client.MCPTool{
		{Name: "tool1", Description: "First tool"},
		{Name: "tool2"},
	}

	p := buildPrompt(config, tools)
	raw := p.Raw

	if raw == "" {
		t.Error("prompt should not be empty")
	}
	// Should mention server name
	if !contains(raw, "test-server") {
		t.Error("prompt should mention server name")
	}
	// Should mention tool names
	if !contains(raw, "tool1") || !contains(raw, "tool2") {
		t.Error("prompt should mention tool names")
	}
}

func TestBuildPrompt_Custom(t *testing.T) {
	config := client.ServerConfig{
		Name:   "my-server",
		Prompt: "You are a helpful assistant.",
	}
	p := buildPrompt(config, nil)

	if p.Raw != "You are a helpful assistant." {
		t.Errorf("custom prompt = %q, want custom prompt", p.Raw)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// buildTestDirective creates an MCPDirective with a mock connection for testing.
func buildTestDirective(conn client.ServerConnection, tools []client.MCPTool) *MCPDirective {
	callables := make([]tool.ToolCallable, len(tools))
	for i, mt := range tools {
		callables[i] = newMCPToolAdapter(mt, conn)
	}

	config := client.ServerConfig{Name: conn.ServerName()}
	p := buildPrompt(config, tools)
	inner := directive.NewStaticDirective(config.Name, p, callables)

	return &MCPDirective{
		inner:      inner,
		connection: conn,
	}
}
