package openai

import "time"

// OpenAIConfig holds the configuration for the OpenAI connector.
type OpenAIConfig struct {
	BaseURL        string
	Model          string
	APIKey         string // Optional
	Timeout        time.Duration
	ResponseFormat *ResponseFormat // Optional: force JSON output mode
}

// DefaultOpenAIConfig creates a configuration with sensible defaults.
func DefaultOpenAIConfig(baseURL, model string) OpenAIConfig {
	return OpenAIConfig{
		BaseURL: baseURL,
		Model:   model,
		APIKey:  "",
		Timeout: 60 * time.Second,
	}
}
