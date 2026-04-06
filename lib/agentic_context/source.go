package agentic_context

// Source represents the origin of a ContextPart or Message.
// It provides type-safe methods to check the source type.
type Source interface {
	IsSystem() bool
	IsUser() bool
	IsAgent() bool
	IsTool() bool
	String() string
}

// systemSource is the concrete implementation for system sources.
type systemSource struct{}

func (s systemSource) IsSystem() bool { return true }
func (s systemSource) IsUser() bool   { return false }
func (s systemSource) IsAgent() bool  { return false }
func (s systemSource) IsTool() bool   { return false }
func (s systemSource) String() string { return string(System) }

// SystemSource is the singleton instance for system sources.
var SystemSource Source = systemSource{}

// userSource is the concrete implementation for user sources.
type userSource struct{}

func (u userSource) IsSystem() bool { return false }
func (u userSource) IsUser() bool   { return true }
func (u userSource) IsAgent() bool  { return false }
func (u userSource) IsTool() bool   { return false }
func (u userSource) String() string { return string(User) }

// UserSource is the singleton instance for user sources.
var UserSource Source = userSource{}

// agentSource is the concrete implementation for agent sources.
type agentSource struct{}

func (a agentSource) IsSystem() bool { return false }
func (a agentSource) IsUser() bool   { return false }
func (a agentSource) IsAgent() bool  { return true }
func (a agentSource) IsTool() bool   { return false }
func (a agentSource) String() string { return string(Agent) }

// AgentSource is the singleton instance for agent sources.
var AgentSource Source = agentSource{}

// toolSource is the concrete implementation for tool sources.
type toolSource struct{}

func (t toolSource) IsSystem() bool { return false }
func (t toolSource) IsUser() bool   { return false }
func (t toolSource) IsAgent() bool  { return false }
func (t toolSource) IsTool() bool   { return true }
func (t toolSource) String() string { return string(Tool) }

// ToolSource is the singleton instance for tool sources.
var ToolSource Source = toolSource{}
