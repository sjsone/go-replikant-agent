package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/connector"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// OpenAIConnector implements the Connector interface for OpenAI-compatible APIs.
type OpenAIConnector struct {
	config OpenAIConfig
	client *http.Client
}

// NewOpenAIConnector creates a new OpenAI connector with the given configuration.
func NewOpenAIConnector(config OpenAIConfig) *OpenAIConnector {
	return &OpenAIConnector{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SendNonStreaming sends messages to the OpenAI API without streaming and returns the complete response content.
// Use this for structured output requests where streaming is not needed.
func (c *OpenAIConnector) SendNonStreaming(ctx context.Context, messages []ChatMessage, responseFormat *ResponseFormat) (string, error) {
	// Build request without streaming
	req := ChatRequest{
		Model:          c.config.Model,
		Messages:       messages,
		Stream:         false,
		ResponseFormat: responseFormat,
	}

	// Marshal request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request with context
	url := c.config.BaseURL + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	// Send request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error.Message != "" {
			return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return "", fmt.Errorf("API error (status %d)", resp.StatusCode)
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract content from first choice
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// SendForRouting implements the Connector interface method for routing decisions.
// It sends a non-streaming request with structured output and returns raw JSON bytes.
func (c *OpenAIConnector) SendForRouting(ctx context.Context, messages []connector.ChatMessage, schema *connector.JSONSchema) (json.RawMessage, error) {
	rf := ResponseFormat{
		Type:       "json_schema",
		JSONSchema: schema,
	}

	content, err := c.SendNonStreaming(ctx, messages, &rf)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(content), nil
}

// Send sends messages to the OpenAI API using streaming and returns the complete response as a ContextPart.
// The response is streamed via onChunk callback while accumulating the complete content.
func (c *OpenAIConnector) Send(ctx context.Context, messages *[]connector.Message, directives []directive.Directive, onChunk connector.ChunkHandler) (error, *agentic_context.ContextPart) {
	// Convert messages to OpenAI format
	openaiMessages := c.messagesToOpenAI(*messages)

	// Build request with streaming enabled
	req := ChatRequest{
		Model:    c.config.Model,
		Messages: openaiMessages,
		Stream:   true,
	}

	// Extract tools from directives and add to request
	allTools := make([]*tool.Tool, 0)
	for _, d := range directives {
		allTools = append(allTools, d.GetTools()...)
	}
	if len(allTools) > 0 {
		req.Tools = c.toolsToOpenAI(allTools)
	} else if c.config.ResponseFormat != nil {
		// Only use response_format when no tools are present.
		// Some servers (e.g. llama.cpp) reject the combination.
		req.ResponseFormat = c.config.ResponseFormat
	}

	// Marshal request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err), nil
	}

	// Create HTTP request with context
	url := c.config.BaseURL + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err), nil
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	// Send request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err), nil
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error.Message != "" {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error.Message), nil
		}
		return fmt.Errorf("API error (status %d)", resp.StatusCode), nil
	}

	var completeContent strings.Builder

	// Accumulate tool calls during streaming - track by index
	type toolCallAccumulator struct {
		call *tool.FunctionCall
		args strings.Builder
	}
	toolCallMap := make(map[int]*toolCallAccumulator)

	// Read SSE stream
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			// Close response body and create partial part with Cancelled=true
			resp.Body.Close()

			// Finalize any accumulated tool calls
			var toolCalls []tool.FunctionCall
			for _, acc := range toolCallMap {
				if acc.call.ID != "" && acc.call.Name != "" {
					acc.call.Arguments = tool.ParseArgs(acc.args.String())
					toolCalls = append(toolCalls, *acc.call)
				}
			}

			// Create partial ContextPart with accumulated content
			part := agentic_context.NewAgentContextPart(completeContent.String())
			part.Cancelled = true

			// Set tool calls if present
			if len(toolCalls) > 0 {
				part.ToolCalls = toolCalls
				part.ToolUse = true
				part.Stop = false
			} else {
				part.Stop = true
			}

			return ctx.Err(), part
		default:
			// Continue processing
		}

		line := scanner.Text()

		// SSE lines start with "data: "
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// "[DONE]" marks end of stream
		if data == "[DONE]" {
			break
		}

		// Parse stream chunk
		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		// Extract content delta and tool calls
		if len(streamResp.Choices) > 0 {
			delta := streamResp.Choices[0].Delta

			// Handle content
			if delta.Content != "" {
				completeContent.WriteString(delta.Content)
				if onChunk != nil {
					onChunk(delta.Content)
				}
			}

			// Handle tool calls - use index to track
			for _, tc := range delta.ToolCalls {
				acc, exists := toolCallMap[tc.Index]
				if !exists {
					// Create new accumulator for this index
					acc = &toolCallAccumulator{
						call: &tool.FunctionCall{
							ID:   tc.ID,
							Name: tc.Function.Name,
						},
					}
					toolCallMap[tc.Index] = acc
				}
				// Update ID if provided (may not be in all deltas)
				if tc.ID != "" {
					acc.call.ID = tc.ID
				}
				// Update name if provided
				if tc.Function.Name != "" {
					acc.call.Name = tc.Function.Name
				}
				// Accumulate arguments
				acc.args.WriteString(tc.Function.Arguments)
			}
		}
	}

	// Finalize all tool calls
	var toolCalls []tool.FunctionCall
	for _, acc := range toolCallMap {
		if acc.call.ID != "" && acc.call.Name != "" {
			acc.call.Arguments = tool.ParseArgs(acc.args.String())
			toolCalls = append(toolCalls, *acc.call)
		}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream read error: %w", err), nil
	}

	// Create ContextPart with complete response
	part := agentic_context.NewAgentContextPart(completeContent.String())

	// Set tool calls if present
	if len(toolCalls) > 0 {
		part.ToolCalls = toolCalls
		part.ToolUse = true
		// Don't stop when there are tool calls - we need to continue to send results back
		part.Stop = false
	} else {
		part.Stop = true
	}

	return nil, part
}

// messagesToOpenAI converts connector messages to OpenAI chat format.
func (c *OpenAIConnector) messagesToOpenAI(messages []connector.Message) []ChatMessage {
	result := make([]ChatMessage, 0, len(messages))
	for _, msg := range messages {
		result = append(result, ChatMessage{
			Role:    c.sourceToRole(msg.Source),
			Content: msg.Text,
		})
	}
	return result
}

// sourceToRole maps a Source to the corresponding OpenAI role.
func (c *OpenAIConnector) sourceToRole(source agentic_context.Source) string {
	switch {
	case source.IsSystem():
		return "system"
	case source.IsUser():
		return "user"
	case source.IsAgent():
		return "assistant"
	case source.IsTool():
		return "tool"
	default:
		return "user"
	}
}

// toolsToOpenAI converts tool.Tool to ToolDefinition format.
func (c *OpenAIConnector) toolsToOpenAI(tools []*tool.Tool) []ToolDefinition {
	result := make([]ToolDefinition, 0, len(tools))
	for _, t := range tools {
		result = append(result, ToolDefinition{
			Type: "function",
			Function: FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}
	return result
}

