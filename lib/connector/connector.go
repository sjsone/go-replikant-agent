package connector

import (
	"context"
	"encoding/json"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/directive"
)

// ChunkHandler is called for each chunk of streaming content from the LLM.
type ChunkHandler func(chunk string)

// The interface which connects to an LLM
type Connector interface {
	// TODO: create SendNonStreaming
	// TODO: maybe rename Send to SendStreaming
	Send(ctx context.Context, messages *[]Message, directives []directive.Directive, onChunk ChunkHandler) (error, *agentic_context.ContextPart)
}

type RoutingConnector interface {
	// SendForRouting sends a request for directive selection/routing.
	// Returns raw JSON bytes; the caller (router) is responsible for parsing.
	SendForRouting(ctx context.Context, messages []ChatMessage) (json.RawMessage, error)
}
