package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/loop"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/router"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// CallRecord records a single call to the connector.
type CallRecord struct {
	Messages   []connector.Message
	Directives []directive.Directive
}

// MockConnector is a mock implementation of the Connector interface for testing.
type MockConnector struct {
	mu                sync.Mutex
	ResponseToReturn  *agentic_context.ContextPart
	ErrorToReturn     error
	RoutingDecision   *router.RoutingDecision
	RoutingError      error
	CallsMade         []CallRecord
	RoutingCallsMade  int
	ChunkHandlerToUse connector.ChunkHandler
}

// NewMockConnector creates a new MockConnector.
func NewMockConnector() *MockConnector {
	return &MockConnector{
		CallsMade: make([]CallRecord, 0),
	}
}

// SetResponse sets the response to return from Send.
func (m *MockConnector) SetResponse(part *agentic_context.ContextPart) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseToReturn = part
}

// SetError sets the error to return from Send.
func (m *MockConnector) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorToReturn = err
}

// SetRoutingDecision sets the routing decision to return from SendForRouting.
func (m *MockConnector) SetRoutingDecision(decision *router.RoutingDecision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RoutingDecision = decision
}

// SetRoutingError sets the error to return from SendForRouting.
func (m *MockConnector) SetRoutingError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RoutingError = err
}

// Send implements the Connector interface.
func (m *MockConnector) Send(ctx context.Context, messages *[]connector.Message, directives []directive.Directive, onChunk connector.ChunkHandler) (error, *agentic_context.ContextPart) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallsMade = append(m.CallsMade, CallRecord{
		Messages:   *messages,
		Directives: directives,
	})

	// Store chunk handler for potential use
	m.ChunkHandlerToUse = onChunk

	// Simulate streaming if we have a response with content
	if m.ResponseToReturn != nil && m.ResponseToReturn.Raw != "" && onChunk != nil {
		onChunk(m.ResponseToReturn.Raw)
	}

	return m.ErrorToReturn, m.ResponseToReturn
}

// SendForRouting implements the RoutingConnector interface.
func (m *MockConnector) SendForRouting(ctx context.Context, messages []connector.ChatMessage, schema *connector.JSONSchema) (json.RawMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RoutingCallsMade++

	if m.RoutingError != nil {
		return nil, m.RoutingError
	}

	if m.RoutingDecision == nil {
		return nil, nil
	}

	raw, err := json.Marshal(m.RoutingDecision)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// GetCallCount returns the number of times Send was called.
func (m *MockConnector) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.CallsMade)
}

// GetLastCall returns the last call made to Send, or nil if no calls were made.
func (m *MockConnector) GetLastCall() *CallRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.CallsMade) == 0 {
		return nil
	}

	last := m.CallsMade[len(m.CallsMade)-1]
	return &last
}

// GetAllCalls returns all calls made to Send.
func (m *MockConnector) GetAllCalls() []CallRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	calls := make([]CallRecord, len(m.CallsMade))
	copy(calls, m.CallsMade)
	return calls
}

// Reset clears all recorded calls and resets responses.
func (m *MockConnector) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallsMade = make([]CallRecord, 0)
	m.RoutingCallsMade = 0
	m.ResponseToReturn = nil
	m.ErrorToReturn = nil
	m.RoutingDecision = nil
	m.RoutingError = nil
	m.ChunkHandlerToUse = nil
}

// ToolExecutionRecord records a single tool execution.
type ToolExecutionRecord struct {
	Call   tool.FunctionCall
	Result tool.FunctionResult
}

// MockDelegate is a mock implementation of SessionDelegate for testing.
type MockDelegate struct {
	mu                 sync.Mutex
	PartsAdded         []*agentic_context.ContextPart
	ToolCallsReceived  [][]tool.FunctionCall
	ToolsExecuted      []ToolExecutionRecord
	ChunksReceived     []string
	RequestsSent       []RequestSentRecord
	LoopIterations     []int
	LoopEnded          bool
	RoutingDecisions   []RoutingDecisionRecord
}

// RequestSentRecord records a request sent event.
type RequestSentRecord struct {
	Messages   []connector.Message
	Directives []directive.Directive
}

// RoutingDecisionRecord records a routing decision event.
type RoutingDecisionRecord struct {
	Decision      router.RoutingDecision
	AllOptions    []*router.RoutingOption
	ActiveOptions []*router.RoutingOption
}

// NewMockDelegate creates a new MockDelegate.
func NewMockDelegate() *MockDelegate {
	return &MockDelegate{
		PartsAdded:         make([]*agentic_context.ContextPart, 0),
		ToolCallsReceived:  make([][]tool.FunctionCall, 0),
		ToolsExecuted:      make([]ToolExecutionRecord, 0),
		ChunksReceived:     make([]string, 0),
		RequestsSent:       make([]RequestSentRecord, 0),
		LoopIterations:     make([]int, 0),
		RoutingDecisions:   make([]RoutingDecisionRecord, 0),
	}
}

// SessionOnPartAdded implements SessionDelegate.
func (m *MockDelegate) SessionOnPartAdded(part *agentic_context.ContextPart) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PartsAdded = append(m.PartsAdded, part)
}

// SessionOnToolCallsReceived implements SessionDelegate.
func (m *MockDelegate) SessionOnToolCallsReceived(calls []tool.FunctionCall) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ToolCallsReceived = append(m.ToolCallsReceived, calls)
}

// SessionOnToolExecuted implements SessionDelegate.
func (m *MockDelegate) SessionOnToolExecuted(call tool.FunctionCall, result tool.FunctionResult) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ToolsExecuted = append(m.ToolsExecuted, ToolExecutionRecord{
		Call:   call,
		Result: result,
	})
}

// SessionOnStreamingChunk implements SessionDelegate.
func (m *MockDelegate) SessionOnStreamingChunk(chunk string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ChunksReceived = append(m.ChunksReceived, chunk)
}

// SessionOnRequestSent implements SessionDelegate.
func (m *MockDelegate) SessionOnRequestSent(messages []connector.Message, directives []directive.Directive) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestsSent = append(m.RequestsSent, RequestSentRecord{
		Messages:   messages,
		Directives: directives,
	})
}

// SessionOnLoopIteration implements SessionDelegate.
func (m *MockDelegate) SessionOnLoopIteration(iteration int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LoopIterations = append(m.LoopIterations, iteration)
}

// SessionOnLoopEnd implements SessionDelegate.
func (m *MockDelegate) SessionOnLoopEnd() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LoopEnded = true
}

// RouterOnRoutingDecision implements router.Delegate.
func (m *MockDelegate) RouterOnRoutingDecision(decision router.RoutingDecision, allOptions []*router.RoutingOption, activeOptions []*router.RoutingOption) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RoutingDecisions = append(m.RoutingDecisions, RoutingDecisionRecord{
		Decision:      decision,
		AllOptions:    allOptions,
		ActiveOptions: activeOptions,
	})
}

// GetPartsAdded returns all parts that were added.
func (m *MockDelegate) GetPartsAdded() []*agentic_context.ContextPart {
	m.mu.Lock()
	defer m.mu.Unlock()

	parts := make([]*agentic_context.ContextPart, len(m.PartsAdded))
	copy(parts, m.PartsAdded)
	return parts
}

// GetToolCallsReceived returns all tool call batches that were received.
func (m *MockDelegate) GetToolCallsReceived() [][]tool.FunctionCall {
	m.mu.Lock()
	defer m.mu.Unlock()

	calls := make([][]tool.FunctionCall, len(m.ToolCallsReceived))
	copy(calls, m.ToolCallsReceived)
	return calls
}

// GetToolsExecuted returns all tool executions that were recorded.
func (m *MockDelegate) GetToolsExecuted() []ToolExecutionRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	executed := make([]ToolExecutionRecord, len(m.ToolsExecuted))
	copy(executed, m.ToolsExecuted)
	return executed
}

// GetChunksReceived returns all chunks that were received.
func (m *MockDelegate) GetChunksReceived() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	chunks := make([]string, len(m.ChunksReceived))
	copy(chunks, m.ChunksReceived)
	return chunks
}

// GetRequestsSent returns all requests that were sent.
func (m *MockDelegate) GetRequestsSent() []RequestSentRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	requests := make([]RequestSentRecord, len(m.RequestsSent))
	copy(requests, m.RequestsSent)
	return requests
}

// GetLoopIterations returns all loop iterations that were recorded.
func (m *MockDelegate) GetLoopIterations() []int {
	m.mu.Lock()
	defer m.mu.Unlock()

	iters := make([]int, len(m.LoopIterations))
	copy(iters, m.LoopIterations)
	return iters
}

// GetRoutingDecisions returns all routing decision events.
func (m *MockDelegate) GetRoutingDecisions() []RoutingDecisionRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	decisions := make([]RoutingDecisionRecord, len(m.RoutingDecisions))
	copy(decisions, m.RoutingDecisions)
	return decisions
}

// Reset clears all recorded events.
func (m *MockDelegate) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.PartsAdded = make([]*agentic_context.ContextPart, 0)
	m.ToolCallsReceived = make([][]tool.FunctionCall, 0)
	m.ToolsExecuted = make([]ToolExecutionRecord, 0)
	m.ChunksReceived = make([]string, 0)
	m.RequestsSent = make([]RequestSentRecord, 0)
	m.LoopIterations = make([]int, 0)
	m.RoutingDecisions = make([]RoutingDecisionRecord, 0)
}

// MockToolCallable is a mock implementation of ToolCallable for testing.
type MockToolCallable struct {
	mu           sync.Mutex
	Name         string
	Response     string
	Error        error
	WasCalled    bool
	ArgsReceived map[string]any
}

// NewMockToolCallable creates a new MockToolCallable.
func NewMockToolCallable(name string, response string) *MockToolCallable {
	return &MockToolCallable{
		Name:         name,
		Response:     response,
		ArgsReceived: make(map[string]any),
	}
}

// NewMockToolCallableWithError creates a MockToolCallable that returns an error.
func NewMockToolCallableWithError(name string, err error) *MockToolCallable {
	return &MockToolCallable{
		Name:         name,
		Error:        err,
		ArgsReceived: make(map[string]any),
	}
}

// Execute implements ToolCallable.
func (m *MockToolCallable) Execute(ctx context.Context, args map[string]any) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.WasCalled = true
	m.ArgsReceived = args

	if m.Error != nil {
		return "", m.Error
	}

	return m.Response, nil
}

// GetTool returns the tool metadata for this callable.
func (m *MockToolCallable) GetTool() *tool.Tool {
	return &tool.Tool{
		Name:        m.Name,
		Description: "Mock tool for testing",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input parameter",
				},
			},
		},
	}
}

// WasCalledWith returns true if the tool was called with the given args.
func (m *MockToolCallable) WasCalledWith(args map[string]any) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.WasCalled {
		return false
	}

	if len(m.ArgsReceived) != len(args) {
		return false
	}

	for k, v := range args {
		if m.ArgsReceived[k] != v {
			return false
		}
	}

	return true
}

// Reset clears the call history.
func (m *MockToolCallable) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.WasCalled = false
	m.ArgsReceived = make(map[string]any)
}

// MockToolCallableWithToolMetadata is a MockToolCallable that embeds tool.Tool.
type MockToolCallableWithToolMetadata struct {
	*MockToolCallable
	tool *tool.Tool
}

// NewMockToolCallableWithToolMetadata creates a MockToolCallable with tool metadata.
func NewMockToolCallableWithToolMetadata(name, description string, response string) *MockToolCallableWithToolMetadata {
	mock := NewMockToolCallable(name, response)
	return &MockToolCallableWithToolMetadata{
		MockToolCallable: mock,
		tool: &tool.Tool{
			Name:        name,
			Description: description,
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"input": map[string]any{
						"type":        "string",
						"description": "Input parameter",
					},
				},
			},
		},
	}
}

// GetTool returns the embedded tool metadata.
func (m *MockToolCallableWithToolMetadata) GetTool() *tool.Tool {
	return m.tool
}

// MockPromptBuilder is a mock PromptBuilder for testing.
type MockPromptBuilder struct {
	prompt     string
	callsMade  int
	directives [][]directive.Directive
}

// NewMockPromptBuilder creates a new MockPromptBuilder.
func NewMockPromptBuilder(prompt string) *MockPromptBuilder {
	return &MockPromptBuilder{
		prompt:     prompt,
		directives: make([][]directive.Directive, 0),
	}
}

// Build implements PromptBuilder.
func (m *MockPromptBuilder) Build(directives []directive.Directive) prompt.Prompt {
	m.callsMade++
	m.directives = append(m.directives, directives)

	return prompt.Prompt{
		Raw: m.prompt,
	}
}

// GetCallCount returns the number of times Build was called.
func (m *MockPromptBuilder) GetCallCount() int {
	return m.callsMade
}

// GetLastDirectives returns the directives from the last call to Build.
func (m *MockPromptBuilder) GetLastDirectives() []directive.Directive {
	if len(m.directives) == 0 {
		return nil
	}
	return m.directives[len(m.directives)-1]
}

// MockLoopController is a mock LoopController for testing.
type MockLoopController struct {
	LoopAgainValue bool
	CallsMade      int
	ContextsSeen   []*agentic_context.AgentContext
}

// NewMockLoopController creates a new MockLoopController.
func NewMockLoopController(loopAgain bool) *MockLoopController {
	return &MockLoopController{
		LoopAgainValue: loopAgain,
		ContextsSeen:   make([]*agentic_context.AgentContext, 0),
	}
}

// LoopAgain implements LoopController.
func (m *MockLoopController) LoopAgain(agentContext *agentic_context.AgentContext) bool {
	m.CallsMade++
	m.ContextsSeen = append(m.ContextsSeen, agentContext)
	return m.LoopAgainValue
}

// GetCallCount returns the number of times LoopAgain was called.
func (m *MockLoopController) GetCallCount() int {
	return m.CallsMade
}

// GetLastContext returns the last context seen by LoopAgain.
func (m *MockLoopController) GetLastContext() *agentic_context.AgentContext {
	if len(m.ContextsSeen) == 0 {
		return nil
	}
	return m.ContextsSeen[len(m.ContextsSeen)-1]
}

// MockMultiplexer is a mock Multiplexer for testing.
type MockMultiplexer struct {
	directives   []directive.Directive
	callsMade    int
	contextsSeen []*agentic_context.AgentContext
}

// NewMockMultiplexer creates a new MockMultiplexer.
func NewMockMultiplexer(directives ...directive.Directive) *MockMultiplexer {
	return &MockMultiplexer{
		directives:   directives,
		contextsSeen: make([]*agentic_context.AgentContext, 0),
	}
}

// GetActiveDirectivesForContext implements Multiplexer.
func (m *MockMultiplexer) GetActiveDirectivesForContext(ctx agentic_context.AgentContext) []directive.Directive {
	m.callsMade++
	m.contextsSeen = append(m.contextsSeen, &ctx)
	return m.directives
}

// SetDirectives sets the directives to return.
func (m *MockMultiplexer) SetDirectives(directives []directive.Directive) {
	m.directives = directives
}

// GetCallCount returns the number of times GetActiveDirectivesForContext was called.
func (m *MockMultiplexer) GetCallCount() int {
	return m.callsMade
}

// GetLastContext returns the last context seen.
func (m *MockMultiplexer) GetLastContext() *agentic_context.AgentContext {
	if len(m.contextsSeen) == 0 {
		return nil
	}
	return m.contextsSeen[len(m.contextsSeen)-1]
}

// AssertMockConnectorReceivedCall asserts that the connector received a specific call.
func AssertMockConnectorReceivedCall(t Testing, connector *MockConnector, expectedDirectives int) {
	t.Helper()
	if connector.GetCallCount() == 0 {
		t.Errorf("Expected connector to receive at least 1 call, got 0")
		return
	}

	lastCall := connector.GetLastCall()
	if lastCall == nil {
		t.Errorf("Expected last call to be non-nil")
		return
	}

	if len(lastCall.Directives) != expectedDirectives {
		t.Errorf("Expected %d directives, got %d", expectedDirectives, len(lastCall.Directives))
	}
}

// AssertMockDelegateReceivedPart asserts that the delegate received a part.
func AssertMockDelegateReceivedPart(t Testing, delegate *MockDelegate) bool {
	t.Helper()
	parts := delegate.GetPartsAdded()
	if len(parts) == 0 {
		t.Errorf("Expected delegate to receive at least 1 part, got 0")
		return false
	}
	return true
}

// AssertMockToolCalled asserts that a mock tool was called.
func AssertMockToolCalled(t Testing, tool *MockToolCallable) bool {
	t.Helper()
	if !tool.WasCalled {
		t.Errorf("Expected tool to be called, but it was not")
		return false
	}
	return true
}

// Testing is an interface for testing.TB.
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// Simple testing implementation
type SimpleTesting struct {
	errors []string
}

func (t *SimpleTesting) Helper() {}

func (t *SimpleTesting) Errorf(format string, args ...any) {
	t.errors = append(t.errors, fmt.Sprintf(format, args...))
}

func (t *SimpleTesting) HasErrors() bool {
	return len(t.errors) > 0
}

func (t *SimpleTesting) GetErrors() []string {
	return t.errors
}

// ErrRoutingFailed is a sentinel error for testing routing failures.
var ErrRoutingFailed = fmt.Errorf("routing failed")

// ErrCommandFailed is a sentinel error for testing command execution failures.
var ErrCommandFailed = fmt.Errorf("command failed")

// MockRouter is a mock implementation of Router for testing.
type MockRouter struct {
	mu                      sync.Mutex
	ResultToReturn          *router.RoutingResult
	ErrorToReturn           error
	RoutingDecisionToReturn *router.RoutingDecision
	CallsMade               int
	Delegate                router.Delegate
}

// NewMockRouter creates a new MockRouter.
func NewMockRouter() *MockRouter {
	return &MockRouter{
		CallsMade: 0,
	}
}

// SetResult sets the routing result to return from Route.
func (m *MockRouter) SetResult(result *router.RoutingResult) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResultToReturn = result
}

// SetRoutingDecision sets the routing decision to return (creates a RoutingResult wrapper).
func (m *MockRouter) SetRoutingDecision(decision *router.RoutingDecision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RoutingDecisionToReturn = decision
}

// SetDelegate sets the delegate for routing decisions.
func (m *MockRouter) SetDelegate(delegate router.Delegate) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Delegate = delegate
}

// SetError causes Route to return nil, simulating a routing error.
func (m *MockRouter) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorToReturn = err
}

// Route implements the Router interface.
func (m *MockRouter) Route(ctx context.Context, userQuery string, allAvailableOptions []*router.RoutingOption) *router.RoutingResult {
	// userQuery parameter is accepted but not used in the mock
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallsMade++

	// If an error is set, return nil to simulate routing failure
	if m.ErrorToReturn != nil {
		return nil
	}

	// If we have a pre-set result, return it
	if m.ResultToReturn != nil {
		result := m.ResultToReturn
		// Call delegate if set
		if m.Delegate != nil && result.Decision != nil {
			m.Delegate.RouterOnRoutingDecision(*result.Decision, allAvailableOptions, result.SelectedOptions)
		}
		return result
	}

	// If we have a routing decision set, create a result from it
	if m.RoutingDecisionToReturn != nil {
		decision := m.RoutingDecisionToReturn

		// Build name map for lookup
		nameMap := make(map[string]*router.RoutingOption, len(allAvailableOptions))
		for _, opt := range allAvailableOptions {
			nameMap[opt.Name] = opt
		}

		selectedOptions := make([]*router.RoutingOption, 0, len(decision.SelectedIDs))
		for _, name := range decision.SelectedIDs {
			if opt, ok := nameMap[name]; ok {
				selectedOptions = append(selectedOptions, opt)
			}
		}

		// Handle fallback behavior
		var finalOptions []*router.RoutingOption
		if len(selectedOptions) == 0 && len(decision.SelectedIDs) == 0 {
			// Empty selection - leave it empty
		} else {
			finalOptions = selectedOptions
		}

		result := &router.RoutingResult{
			SelectedOptions: finalOptions,
			Decision:        decision,
		}
		// Call delegate if set
		if m.Delegate != nil {
			m.Delegate.RouterOnRoutingDecision(*decision, allAvailableOptions, finalOptions)
		}
		return result
	}

	// Default: return all options with a nil decision
	return &router.RoutingResult{
		SelectedOptions: allAvailableOptions,
		Decision:        nil,
	}
}

// GetCallCount returns the number of times Route was called.
func (m *MockRouter) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.CallsMade
}

// Reset clears the call history and resets responses.
func (m *MockRouter) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallsMade = 0
	m.ResultToReturn = nil
	m.RoutingDecisionToReturn = nil
	m.Delegate = nil
}

// MockLoopControllerDelegate is a mock implementation of loop.Delegate for testing.
type MockLoopControllerDelegate struct {
	mu        sync.Mutex
	Decisions []loop.LoopDecision
}

// NewMockLoopControllerDelegate creates a new MockLoopControllerDelegate.
func NewMockLoopControllerDelegate() *MockLoopControllerDelegate {
	return &MockLoopControllerDelegate{
		Decisions: make([]loop.LoopDecision, 0),
	}
}

// LoopOnLoopDecision implements loop.Delegate.
func (m *MockLoopControllerDelegate) LoopOnLoopDecision(decision loop.LoopDecision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Decisions = append(m.Decisions, decision)
}

// GetDecisions returns all recorded decisions.
func (m *MockLoopControllerDelegate) GetDecisions() []loop.LoopDecision {
	m.mu.Lock()
	defer m.mu.Unlock()

	decisions := make([]loop.LoopDecision, len(m.Decisions))
	copy(decisions, m.Decisions)
	return decisions
}

// GetLastDecision returns the last recorded decision, or nil if none.
func (m *MockLoopControllerDelegate) GetLastDecision() *loop.LoopDecision {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.Decisions) == 0 {
		return nil
	}
	d := m.Decisions[len(m.Decisions)-1]
	return &d
}

// GetDecisionCount returns the number of recorded decisions.
func (m *MockLoopControllerDelegate) GetDecisionCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Decisions)
}

// Reset clears all recorded decisions.
func (m *MockLoopControllerDelegate) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Decisions = make([]loop.LoopDecision, 0)
}
