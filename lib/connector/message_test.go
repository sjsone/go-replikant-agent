package connector

import (
	"testing"

	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
)

func TestNewSystemMessage(t *testing.T) {
	text := "System instructions"
	msg := NewSystemMessage(text)

	if msg.Text != text {
		t.Errorf("Text = %q, want %q", msg.Text, text)
	}

	if !msg.IsSystemMessage() {
		t.Error("IsSystemMessage() returned false")
	}

	if msg.IsUserMessage() {
		t.Error("IsUserMessage() returned true")
	}

	if msg.IsAgentMessage() {
		t.Error("IsAgentMessage() returned true")
	}

	if msg.Source != agentic_context.SystemSource {
		t.Errorf("Source = %v, want SystemSource", msg.Source)
	}
}

func TestNewUserMessage(t *testing.T) {
	text := "User input"
	msg := NewUserMessage(text)

	if msg.Text != text {
		t.Errorf("Text = %q, want %q", msg.Text, text)
	}

	if !msg.IsUserMessage() {
		t.Error("IsUserMessage() returned false")
	}

	if msg.IsSystemMessage() {
		t.Error("IsSystemMessage() returned true")
	}

	if msg.IsAgentMessage() {
		t.Error("IsAgentMessage() returned true")
	}

	if msg.Source != agentic_context.UserSource {
		t.Errorf("Source = %v, want UserSource", msg.Source)
	}
}

func TestNewAgentMessage(t *testing.T) {
	text := "Agent response"
	msg := NewAgentMessage(text)

	if msg.Text != text {
		t.Errorf("Text = %q, want %q", msg.Text, text)
	}

	if !msg.IsAgentMessage() {
		t.Error("IsAgentMessage() returned false")
	}

	if msg.IsSystemMessage() {
		t.Error("IsSystemMessage() returned true")
	}

	if msg.IsUserMessage() {
		t.Error("IsUserMessage() returned true")
	}

	if msg.Source != agentic_context.AgentSource {
		t.Errorf("Source = %v, want AgentSource", msg.Source)
	}
}

func TestMessageIsXMethods(t *testing.T) {
	tests := []struct {
		name     string
		message  Message
		isSystem bool
		isUser   bool
		isAgent  bool
	}{
		{
			name:     "system message",
			message:  NewSystemMessage("System text"),
			isSystem: true,
			isUser:   false,
			isAgent:  false,
		},
		{
			name:     "user message",
			message:  NewUserMessage("User text"),
			isSystem: false,
			isUser:   true,
			isAgent:  false,
		},
		{
			name:     "agent message",
			message:  NewAgentMessage("Agent text"),
			isSystem: false,
			isUser:   false,
			isAgent:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.message.IsSystemMessage() != tt.isSystem {
				t.Errorf("IsSystemMessage() = %v, want %v", tt.message.IsSystemMessage(), tt.isSystem)
			}

			if tt.message.IsUserMessage() != tt.isUser {
				t.Errorf("IsUserMessage() = %v, want %v", tt.message.IsUserMessage(), tt.isUser)
			}

			if tt.message.IsAgentMessage() != tt.isAgent {
				t.Errorf("IsAgentMessage() = %v, want %v", tt.message.IsAgentMessage(), tt.isAgent)
			}
		})
	}
}

func TestMessage_EmptyText(t *testing.T) {
	tests := []struct {
		name    string
		message Message
	}{
		{
			name:    "empty system message",
			message: NewSystemMessage(""),
		},
		{
			name:    "empty user message",
			message: NewUserMessage(""),
		},
		{
			name:    "empty agent message",
			message: NewAgentMessage(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.message.Text != "" {
				t.Errorf("Text = %q, want empty string", tt.message.Text)
			}
		})
	}
}

func TestMessage_MultilineText(t *testing.T) {
	text := "Line 1\nLine 2\nLine 3"
	msg := NewUserMessage(text)

	if msg.Text != text {
		t.Errorf("Text = %q, want %q", msg.Text, text)
	}
}

func TestMessage_UnicodeText(t *testing.T) {
	text := "Hello 世界 🚀"
	msg := NewUserMessage(text)

	if msg.Text != text {
		t.Errorf("Text = %q, want %q", msg.Text, text)
	}
}

func TestMessage_StructFields(t *testing.T) {
	text := "Test message"
	msg := NewSystemMessage(text)

	// Test that we can access the fields directly
	if msg.Source == nil {
		t.Error("Source is nil")
	}

	if msg.Text == "" {
		t.Error("Text is empty")
	}
}

func TestMessage_ZeroValue(t *testing.T) {
	var msg Message

	// Zero value message should have nil source
	if msg.Source != nil {
		t.Errorf("Zero value Source = %v, want nil", msg.Source)
	}

	// Zero value message should have empty text
	if msg.Text != "" {
		t.Errorf("Zero value Text = %q, want empty string", msg.Text)
	}

	// Type check methods will panic with nil source
	// This is expected behavior - messages should always be created
	// using the constructor functions
	t.Run("nil source panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when calling IsSystemMessage() with nil source")
			}
		}()
		_ = msg.IsSystemMessage()
	})
}

func TestMessage_LongText(t *testing.T) {
	// Create a long text
	longText := ""
	for i := 0; i < 1000; i++ {
		longText += "word "
	}

	msg := NewUserMessage(longText)

	if msg.Text != longText {
		t.Errorf("Long text was not preserved correctly")
	}

	if len(msg.Text) != len(longText) {
		t.Errorf("Long text length = %d, want %d", len(msg.Text), len(longText))
	}
}

func TestMessage_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "tabs",
			text: "Text\twith\ttabs",
		},
		{
			name: "newlines",
			text: "Text\nwith\nnewlines",
		},
		{
			name: "quotes",
			text: `Text with "quotes" and 'apostrophes'`,
		},
		{
			name: "backslashes",
			text: "Text\\with\\backslashes",
		},
		{
			name: "mixed special chars",
			text: "Text\n\twith\"mixed'\\special\\chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewUserMessage(tt.text)
			if msg.Text != tt.text {
				t.Errorf("Text = %q, want %q", msg.Text, tt.text)
			}
		})
	}
}

func TestMessage_SourceString(t *testing.T) {
	tests := []struct {
		name        string
		message     Message
		expectedStr string
	}{
		{
			name:        "system message",
			message:     NewSystemMessage("test"),
			expectedStr: string(agentic_context.System),
		},
		{
			name:        "user message",
			message:     NewUserMessage("test"),
			expectedStr: string(agentic_context.User),
		},
		{
			name:        "agent message",
			message:     NewAgentMessage("test"),
			expectedStr: string(agentic_context.Agent),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.message.Source.String() != tt.expectedStr {
				t.Errorf("Source.String() = %q, want %q", tt.message.Source.String(), tt.expectedStr)
			}
		})
	}
}

func TestMessage_ImmutableSource(t *testing.T) {
	// Verify that changing the source reference doesn't affect the message
	msg := NewSystemMessage("test")

	// Get the source
	originalSource := msg.Source

	// The source should be a singleton, so changing the reference
	// doesn't affect the original
	if originalSource != agentic_context.SystemSource {
		t.Error("Source is not the SystemSource singleton")
	}
}

func TestMessage_Copy(t *testing.T) {
	original := NewUserMessage("Original text")

	// Create a copy by copying the struct
	copy := original

	if copy.Text != original.Text {
		t.Errorf("Copy Text = %q, want %q", copy.Text, original.Text)
	}

	if copy.Source != original.Source {
		t.Errorf("Copy Source = %v, want %v", copy.Source, original.Source)
	}

	// Modify the copy
	copy.Text = "Modified text"

	// Original should be unchanged
	if original.Text != "Original text" {
		t.Errorf("Original Text = %q, want 'Original text'", original.Text)
	}
}
