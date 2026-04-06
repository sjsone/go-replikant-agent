package agentic_context

import (
	"testing"
)

func TestSystemSource(t *testing.T) {
	tests := []struct {
		name     string
		source   Source
		isSystem bool
		isUser   bool
		isAgent  bool
		isTool   bool
		str      string
	}{
		{
			name:     "SystemSource",
			source:   SystemSource,
			isSystem: true,
			isUser:   false,
			isAgent:  false,
			isTool:   false,
			str:      string(System),
		},
		{
			name:     "UserSource",
			source:   UserSource,
			isSystem: false,
			isUser:   true,
			isAgent:  false,
			isTool:   false,
			str:      string(User),
		},
		{
			name:     "AgentSource",
			source:   AgentSource,
			isSystem: false,
			isUser:   false,
			isAgent:  true,
			isTool:   false,
			str:      string(Agent),
		},
		{
			name:     "ToolSource",
			source:   ToolSource,
			isSystem: false,
			isUser:   false,
			isAgent:  false,
			isTool:   true,
			str:      string(Tool),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.source.IsSystem() != tt.isSystem {
				t.Errorf("IsSystem() = %v, want %v", tt.source.IsSystem(), tt.isSystem)
			}

			if tt.source.IsUser() != tt.isUser {
				t.Errorf("IsUser() = %v, want %v", tt.source.IsUser(), tt.isUser)
			}

			if tt.source.IsAgent() != tt.isAgent {
				t.Errorf("IsAgent() = %v, want %v", tt.source.IsAgent(), tt.isAgent)
			}

			if tt.source.IsTool() != tt.isTool {
				t.Errorf("IsTool() = %v, want %v", tt.source.IsTool(), tt.isTool)
			}

			if tt.source.String() != tt.str {
				t.Errorf("String() = %v, want %v", tt.source.String(), tt.str)
			}
		})
	}
}

func TestSourceSingletons(t *testing.T) {
	// Verify that the sources are singletons (same instance)
	if SystemSource != SystemSource {
		t.Error("SystemSource is not a singleton")
	}

	if UserSource != UserSource {
		t.Error("UserSource is not a singleton")
	}

	if AgentSource != AgentSource {
		t.Error("AgentSource is not a singleton")
	}

	if ToolSource != ToolSource {
		t.Error("ToolSource is not a singleton")
	}

	// Verify they're all different from each other
	sources := []Source{SystemSource, UserSource, AgentSource, ToolSource}
	for i := 0; i < len(sources); i++ {
		for j := i + 1; j < len(sources); j++ {
			if sources[i] == sources[j] {
				t.Errorf("Sources %d and %d are the same instance", i, j)
			}
		}
	}
}

func TestSourceTypes(t *testing.T) {
	t.Run("SystemSource type checks", func(t *testing.T) {
		if !SystemSource.IsSystem() {
			t.Error("SystemSource.IsSystem() should be true")
		}
		if SystemSource.IsUser() {
			t.Error("SystemSource.IsUser() should be false")
		}
		if SystemSource.IsAgent() {
			t.Error("SystemSource.IsAgent() should be false")
		}
		if SystemSource.IsTool() {
			t.Error("SystemSource.IsTool() should be false")
		}
		if SystemSource.String() != string(System) {
			t.Errorf("SystemSource.String() = %q, want %q", SystemSource.String(), string(System))
		}
	})

	t.Run("UserSource type checks", func(t *testing.T) {
		if UserSource.IsSystem() {
			t.Error("UserSource.IsSystem() should be false")
		}
		if !UserSource.IsUser() {
			t.Error("UserSource.IsUser() should be true")
		}
		if UserSource.IsAgent() {
			t.Error("UserSource.IsAgent() should be false")
		}
		if UserSource.IsTool() {
			t.Error("UserSource.IsTool() should be false")
		}
		if UserSource.String() != string(User) {
			t.Errorf("UserSource.String() = %q, want %q", UserSource.String(), string(User))
		}
	})

	t.Run("AgentSource type checks", func(t *testing.T) {
		if AgentSource.IsSystem() {
			t.Error("AgentSource.IsSystem() should be false")
		}
		if AgentSource.IsUser() {
			t.Error("AgentSource.IsUser() should be false")
		}
		if !AgentSource.IsAgent() {
			t.Error("AgentSource.IsAgent() should be true")
		}
		if AgentSource.IsTool() {
			t.Error("AgentSource.IsTool() should be false")
		}
		if AgentSource.String() != string(Agent) {
			t.Errorf("AgentSource.String() = %q, want %q", AgentSource.String(), string(Agent))
		}
	})

	t.Run("ToolSource type checks", func(t *testing.T) {
		if ToolSource.IsSystem() {
			t.Error("ToolSource.IsSystem() should be false")
		}
		if ToolSource.IsUser() {
			t.Error("ToolSource.IsUser() should be false")
		}
		if ToolSource.IsAgent() {
			t.Error("ToolSource.IsAgent() should be false")
		}
		if !ToolSource.IsTool() {
			t.Error("ToolSource.IsTool() should be true")
		}
		if ToolSource.String() != string(Tool) {
			t.Errorf("ToolSource.String() = %q, want %q", ToolSource.String(), string(Tool))
		}
	})
}

func TestSourceStringConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"System constant", string(System)},
		{"User constant", string(User)},
		{"Agent constant", string(Agent)},
		{"Tool constant", string(Tool)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Expected %s to be non-empty", tt.name)
			}
		})
	}
}

func TestSourceInterfaceImplementation(t *testing.T) {
	// Verify all sources implement the Source interface
	var _ Source = SystemSource
	var _ Source = UserSource
	var _ Source = AgentSource
	var _ Source = ToolSource
}
