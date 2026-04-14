package simple

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/router"
)

// mockRoutingConnector is a test double for RoutingConnector.
type mockRoutingConnector struct {
	decision *router.RoutingDecision
	err      error
}

func (m *mockRoutingConnector) SendForRouting(ctx context.Context, messages []connector.Message, schema *connector.JSONSchema) (json.RawMessage, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.decision == nil {
		return nil, nil
	}
	raw, err := json.Marshal(m.decision)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func TestSimpleRouter_Route_Success(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation capabilities"},
		{Name: "Weather", Text: "Weather-related tools"},
		{Name: "Stocks", Text: "Stock market tools"},
	}

	expectedDecision := &router.RoutingDecision{
		SelectedIDs: []string{"General", "Weather"},
		Reasoning:   "User wants general conversation and weather info",
	}

	mockConn := &mockRoutingConnector{
		decision: expectedDecision,
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", options)

	// Verify
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}

	result := routingResult.SelectedOptions
	if len(result) != 2 {
		t.Errorf("Expected 2 selected options, got %d", len(result))
	}

	if result[0].Name != "General" {
		t.Errorf("Expected first option to be 'General', got '%s'", result[0].Name)
	}

	if result[1].Name != "Weather" {
		t.Errorf("Expected second option to be 'Weather', got '%s'", result[1].Name)
	}

	if routingResult.Decision == nil {
		t.Error("Expected decision to be non-nil")
	}

	if routingResult.Decision.SelectedIDs[0] != "General" || routingResult.Decision.SelectedIDs[1] != "Weather" {
		t.Errorf("Expected decision IDs to be [General, Weather], got %v", routingResult.Decision.SelectedIDs)
	}
}

func TestSimpleRouter_Route_EmptyOptions(t *testing.T) {
	// Setup
	mockConn := &mockRoutingConnector{
		decision: &router.RoutingDecision{SelectedIDs: []string{}},
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", []*router.RoutingOption{})

	// Verify
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}

	result := routingResult.SelectedOptions
	if len(result) != 0 {
		t.Errorf("Expected 0 options for empty input, got %d", len(result))
	}

	if routingResult.Decision != nil {
		t.Error("Expected decision to be nil for empty options")
	}
}

func TestSimpleRouter_Route_ConnectorError_FallbackNone(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation"},
		{Name: "Weather", Text: "Weather tools"},
	}

	mockConn := &mockRoutingConnector{
		err: &testError{"routing failed"},
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", options)

	// Verify - returns error on connector error
	if err == nil {
		t.Errorf("Expected error on connector error, got nil")
	}
	if routingResult != nil {
		t.Errorf("Expected nil result on connector error, got non-nil")
	}
}

func TestSimpleRouter_Route_ConnectorError_FallbackAll(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation"},
		{Name: "Weather", Text: "Weather tools"},
	}

	mockConn := &mockRoutingConnector{
		err: &testError{"routing failed"},
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", options)

	// Verify - returns error on connector error
	if err == nil {
		t.Errorf("Expected error on connector error, got nil")
	}
	if routingResult != nil {
		t.Errorf("Expected nil result on connector error, got non-nil")
	}
}

func TestSimpleRouter_Route_ConnectorError_FallbackError(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation"},
		{Name: "Weather", Text: "Weather tools"},
	}

	mockConn := &mockRoutingConnector{
		err: &testError{"routing failed"},
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", options)

	// Verify - returns error on connector error
	if err == nil {
		t.Errorf("Expected error on connector error, got nil")
	}
	if routingResult != nil {
		t.Errorf("Expected nil result on connector error, got non-nil")
	}
}

func TestSimpleRouter_Route_EmptyIDs(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation"},
		{Name: "Weather", Text: "Weather tools"},
	}

	expectedDecision := &router.RoutingDecision{
		SelectedIDs: []string{}, // Empty selection
		Reasoning:   "No options are relevant",
	}

	mockConn := &mockRoutingConnector{
		decision: expectedDecision,
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", options)

	// Verify
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}

	result := routingResult.SelectedOptions
	if len(result) != 0 {
		t.Errorf("Expected 0 options for empty IDs, got %d", len(result))
	}

	if routingResult.Decision == nil {
		t.Error("Expected decision to be non-nil")
	}

	if len(routingResult.Decision.SelectedIDs) != 0 {
		t.Errorf("Expected empty decision IDs, got %v", routingResult.Decision.SelectedIDs)
	}
}

func TestSimpleRouter_SetExampleMessages(t *testing.T) {
	// Setup
	examples := []connector.Message{
		connector.NewUserMessage("Example user message"),
		connector.NewAgentMessage(`{"selected_ids": ["General"], "reasoning": "Example"}`),
	}

	mockConn := &mockRoutingConnector{
		decision: &router.RoutingDecision{SelectedIDs: []string{"General"}},
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)
	rtr.SetExampleMessages(examples)

	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation"},
	}

	// Execute
	ctx := context.Background()
	routingResult, err := rtr.Route(ctx, "test user query", options)

	// Verify
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}

	result := routingResult.SelectedOptions
	if len(result) != 1 {
		t.Errorf("Expected 1 selected option, got %d", len(result))
	}

	if routingResult.Decision == nil {
		t.Error("Expected decision to be non-nil")
	}
}

func TestSimpleRouter_SetDelegate(t *testing.T) {
	// Setup
	mockConn := &mockRoutingConnector{
		decision: &router.RoutingDecision{SelectedIDs: []string{}},
	}

	rtr := NewSimpleRouter("Test routing prompt", mockConn)

	// Setting nil delegate should not panic
	rtr.SetDelegate(nil)

	// If we got here without panicking, the test passes
}

func TestSimpleRouter_BuildSystemPrompt(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "General", Text: "General conversation"},
		{Name: "Weather", Text: "Weather information"},
	}

	rtr := NewSimpleRouter("Custom prompt: ", nil)

	// Execute
	prompt := rtr.buildSystemPrompt(options)

	// Verify
	if prompt == "" {
		t.Error("Expected non-empty system prompt")
	}

	// Check that option names are included
	if len(prompt) < len("General") || len(prompt) < len("Weather") {
		t.Error("Expected system prompt to include option names")
	}
}

func TestSimpleRouter_FilterOptionsByName(t *testing.T) {
	// Setup
	options := []*router.RoutingOption{
		{Name: "Option0", Text: "Zero"},
		{Name: "Option1", Text: "One"},
		{Name: "Option2", Text: "Two"},
	}

	// Execute - select names "Option0" and "Option2"
	result, err := filterOptionsByName(options, []string{"Option0", "Option2"})

	// Verify
	if err != nil {
		t.Fatalf("filterOptionsByName returned error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 filtered options, got %d", len(result))
	}

	if result[0].Name != "Option0" {
		t.Errorf("Expected first option to be 'Option0', got '%s'", result[0].Name)
	}

	if result[1].Name != "Option2" {
		t.Errorf("Expected second option to be 'Option2', got '%s'", result[1].Name)
	}
}

// Test helper types

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
