package ai

// ChatRequest is the HTTP request body for POST /v1/ai/chat.
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse is the full (non-streaming) response.
type ChatResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Provider string `json:"provider"`
	Usage    *Usage `json:"usage,omitempty"`
}

// Usage tracks token consumption for a single request.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// StreamEvent is a single SSE chunk sent during streaming.
type StreamEvent struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
	Usage   *Usage `json:"usage,omitempty"`
}

// ProviderRequest is the normalized request sent to any LLM provider.
type ProviderRequest struct {
	Model       string    `json:"model"`
	System      string    `json:"system"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// Message is a single chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ProviderResponse is the normalized response from any LLM provider.
type ProviderResponse struct {
	Content string
	Usage   *Usage
}

// StreamChunk is a normalized streaming chunk from any provider.
type StreamChunk struct {
	Content string
	Done    bool
	Usage   *Usage
	Err     error
}
