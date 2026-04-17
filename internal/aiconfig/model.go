package aiconfig

import "time"

type Provider struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	ProviderName  string    `json:"provider"`
	ModelID       string    `json:"model_id"`
	HasAPIKey     bool      `json:"has_api_key"`
	BaseURL       *string   `json:"base_url,omitempty"`
	MaxTokens     int       `json:"max_tokens"`
	Temperature   float64   `json:"temperature"`
	MonthlyUsage  int       `json:"monthly_usage"`
	MonthlyLimit  int       `json:"monthly_limit"`
	LastResetAt   time.Time `json:"last_reset_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateRequest struct {
	Provider    string  `json:"provider"`
	ModelID     string  `json:"model_id"`
	APIKey      string  `json:"api_key,omitempty"`
	BaseURL     *string `json:"base_url,omitempty"`
	MaxTokens   *int    `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
}

type UpdateRequest struct {
	Provider    *string  `json:"provider,omitempty"`
	ModelID     *string  `json:"model_id,omitempty"`
	APIKey      *string  `json:"api_key,omitempty"`
	BaseURL     *string  `json:"base_url,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
}

var ValidProviders = map[string]bool{
	"anthropic": true,
	"openai":    true,
	"ollama":    true,
}
