package connector

import "github.com/sjsone/go-replikant-agent/lib/agentic_context"

type Message struct {
	Source agentic_context.Source
	Text   string
}

func (m *Message) Role() string {
	switch {
	case m.Source.IsSystem():
		return "system"
	case m.Source.IsUser():
		return "user"
	case m.Source.IsAgent():
		return "assistant"
	case m.Source.IsTool():
		return "tool"
	default:
		return "user"
	}
}

func NewSystemMessage(text string) Message {
	return Message{
		Source: agentic_context.SystemSource,
		Text:   text,
	}
}

func NewUserMessage(text string) Message {
	return Message{
		Source: agentic_context.UserSource,
		Text:   text,
	}
}

func NewAgentMessage(text string) Message {
	return Message{
		Source: agentic_context.AgentSource,
		Text:   text,
	}
}

func (m *Message) IsSystemMessage() bool {
	return m.Source.IsSystem()
}

func (m *Message) IsUserMessage() bool {
	return m.Source.IsUser()
}

func (m *Message) IsAgentMessage() bool {
	return m.Source.IsAgent()
}
