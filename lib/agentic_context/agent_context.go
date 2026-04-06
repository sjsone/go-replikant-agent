package agentic_context

type AgentContext struct {
	Parts []*ContextPart
}

func NewAgentContext() *AgentContext {
	return &AgentContext{
		Parts: make([]*ContextPart, 0),
	}
}

func (c *AgentContext) AddPart(p *ContextPart) {
	c.Parts = append(c.Parts, p)
}

func (c *AgentContext) GetLatestPart() *ContextPart {
	return c.Parts[len(c.Parts)-1]
}

// GetLatestNonToolPart walks Parts backwards and returns the first part
// where ToolUse == false and IsToolResult() == false. Returns nil if no such part exists.
func (c *AgentContext) GetLatestNonToolPart() *ContextPart {
	for i := len(c.Parts) - 1; i >= 0; i-- {
		p := c.Parts[i]
		if !p.ToolUse && !p.IsToolResult() {
			return p
		}
	}
	return nil
}
