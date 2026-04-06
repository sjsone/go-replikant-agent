package simple

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/router"
)

type SimpleRouter struct {
	Prompt          string // TODO: check if this should be of type `prompt.Prompt` instead of `string`
	delegate        router.Delegate
	connector       connector.RoutingConnector
	exampleMessages []connector.ChatMessage
	userQuery       string
}

func NewSimpleRouter(
	prompt string,
	conn connector.RoutingConnector,
) *SimpleRouter {
	r := &SimpleRouter{
		Prompt:          prompt,
		connector:       conn,
		exampleMessages: make([]connector.ChatMessage, 0),
	}
	// Initialize with default few-shot examples
	r.exampleMessages = r.buildExampleMessages()
	return r
}

func (r *SimpleRouter) SetExampleMessages(exampleMessages []connector.ChatMessage) {
	r.exampleMessages = exampleMessages
}

func (r *SimpleRouter) SetDelegate(delegate router.Delegate) {
	r.delegate = delegate
}

func (r *SimpleRouter) Route(ctx context.Context, userQuery string, allAvailableOptions []*router.RoutingOption) *router.RoutingResult {
	// Store the user query for use in buildUserPrompt
	r.userQuery = userQuery

	// Handle empty options case
	if len(allAvailableOptions) == 0 {
		return &router.RoutingResult{
			SelectedOptions: []*router.RoutingOption{},
			Decision:        nil,
		}
	}

	// Build routing messages
	messages := r.buildRoutingMessages(allAvailableOptions)

	// Call the connector
	raw, err := r.connector.SendForRouting(ctx, messages)
	if err != nil {
		return nil
	}

	// Parse the raw JSON into a RoutingDecision
	decision, err := parseRoutingDecision(raw)
	if err != nil {
		return nil
	}

	// Filter options based on decision
	result := filterOptionsByName(allAvailableOptions, decision.SelectedIDs)

	if r.delegate != nil {
		r.delegate.RouterOnRoutingDecision(*decision, allAvailableOptions, result)
	}

	return &router.RoutingResult{
		SelectedOptions: result,
		Decision:        decision,
	}
}

// parseRoutingDecision parses raw JSON into a RoutingDecision.
func parseRoutingDecision(raw json.RawMessage) (*router.RoutingDecision, error) {
	var rawDecision struct {
		ActiveIndices   []int    `json:"active_indices"`
		SelectedIndices []int    `json:"selected_indices"`
		SelectedIDs     []string `json:"selected_ids"`
		Reasoning       string   `json:"reasoning"`
		Confidence      float64  `json:"confidence"`
	}
	if err := json.Unmarshal(raw, &rawDecision); err != nil {
		return nil, fmt.Errorf("failed to parse routing decision: %w", err)
	}

	// Use selected_ids if provided (new format)
	ids := rawDecision.SelectedIDs

	return &router.RoutingDecision{
		SelectedIDs: ids,
		Reasoning:   rawDecision.Reasoning,
		Confidence:  rawDecision.Confidence,
	}, nil
}

// mapIndicesToNames converts integer indices to directive names using the provided map.
// Unknown indices are silently skipped.
func mapIndicesToNames(indices []int, availableNames map[int]string) []string {
	names := make([]string, 0, len(indices))
	for _, idx := range indices {
		if name, ok := availableNames[idx]; ok {
			names = append(names, name)
		}
	}
	return names
}

// buildRoutingMessages builds the complete message sequence including examples and user prompt.
func (r *SimpleRouter) buildRoutingMessages(allAvailableOptions []*router.RoutingOption) []connector.ChatMessage {
	systemPrompt := r.buildSystemPrompt(allAvailableOptions)

	messages := []connector.ChatMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Add example messages if provided
	messages = append(messages, r.exampleMessages...)

	// Add the user prompt
	messages = append(messages, connector.ChatMessage{
		Role:    "user",
		Content: r.buildUserPrompt(),
	})

	return messages
}

// buildSystemPrompt constructs the system prompt describing available options.
func (r *SimpleRouter) buildSystemPrompt(options []*router.RoutingOption) string {
	prompt := r.Prompt
	if prompt == "" {
		prompt += "You are a routing assistant. Select which options are relevant based on the user's request.\n"
		prompt += "\n"
		prompt += "CRITICAL RULE: When analyzing tool dependencies, you MUST put ALL relevant directive names in the selected_ids array.\n"
		prompt += "CRITICAL RULE: You MUST use bullet-point-style sentences in the reasoning. \n"
		prompt += "\n"
		prompt += "Tool Dependency Pattern: If tool A requires data that tool B provides, then BOTH directives must be activated.\n"
		prompt += "\n"
	}

	prompt += "Available options:\n"
	for _, opt := range options {
		prompt += fmt.Sprintf("[%s] %s\n", opt.Name, opt.Text)
	}
	prompt += "\nReturn JSON with selected_ids array containing ALL names of relevant options."

	return prompt
}

func (r *SimpleRouter) buildExampleMessages() []connector.ChatMessage {
	examples := []connector.ChatMessage{
		// Example 1: Single directive - simple capability request
		{
			Role:    "user",
			Content: "Calculate the sum of these numbers",
		},
		{
			Role:    "assistant",
			Content: `{"selected_ids": ["calculator"], "reasoning": "User needs mathematical calculation. The calculator directive provides calculation tools."}`,
		},
		// Example 2: Multiple directives with tool dependency
		{
			Role:    "user",
			Content: "Generate a chart for the sales data",
		},
		{
			Role:    "assistant",
			Content: `{"selected_ids": ["data_retrieval", "chart_generation"], "reasoning": "User wants a chart generated from sales data. Directive data_retrieval has data_retrieval tool. Directive chart_generation has chart_generation tool which PROCESSES the retrieved data into visual format. Tool dependency: chart_generation needs input from data_retrieval. Therefore BOTH directives must be activated."}`,
		},
		// Example 3: No directives - general conversation
		{
			Role:    "user",
			Content: "Hello, can you help me?",
		},
		{
			Role:    "assistant",
			Content: `{"selected_ids": [], "reasoning": "General conversational opening. Doesn't require any tools or capabilities. No directives needed."}`,
		},
		// Example 4: Multiple dependent directives for complex workflow
		{
			Role:    "user",
			Content: "Translate this document and format it as a PDF report",
		},
		{
			Role:    "assistant",
			Content: `{"selected_ids": ["translation", "document_formatting", "pdf_generation"], "reasoning": "User requests translation with PDF formatting. Directive translation provides translation capability. Directive document_formatting provides document formatting tools. Directive pdf_generation provides PDF generation. All three directives are needed: translation for content conversion, formatting for layout, and PDF generation for final output format."}`,
		},
		// Example 5: Optional directive selection based on user preference
		{
			Role:    "user",
			Content: "Send me a summary of today's news",
		},
		{
			Role:    "assistant",
			Content: `{"selected_ids": ["news"], "reasoning": "User wants a news summary. Directive news provides news retrieval and summarization tools. This is sufficient for the request. No additional directives are required."}`,
		},
	}
	return examples
}

// buildUserPrompt constructs the user prompt for the routing request.
func (r *SimpleRouter) buildUserPrompt() string {
	return fmt.Sprintf("User query: %s\n\nWhich options should be selected? Return JSON with selected_ids array containing ALL relevant option names.", r.userQuery)
}

// filterOptionsByName returns options matching the given names.
// An empty names array means no options should be selected.
func filterOptionsByName(options []*router.RoutingOption, names []string) []*router.RoutingOption {
	nameMap := make(map[string]*router.RoutingOption, len(options))
	for _, opt := range options {
		nameMap[opt.Name] = opt
	}

	result := make([]*router.RoutingOption, 0, len(names))
	for _, name := range names {
		if opt, ok := nameMap[name]; ok {
			result = append(result, opt)
		}
	}
	return result
}
