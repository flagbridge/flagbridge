// Package provider defines the abstraction for LLM providers used by the AI layer
// per ADR-0006 (docs/adr/0006-ai-provider-interface.md).
//
// The Provider interface is intentionally minimal — a single Complete method
// returning a channel of Event values — so adding a new provider (Gemini, Mistral,
// Bedrock, Ollama) never requires changes to handlers, middleware, or rate limiters.
//
// Events are emitted in order:
//
//  1. Zero or more Event{Delta: "..."} as tokens stream
//  2. Exactly one terminating Event, one of:
//     - Event{Usage: ...} on successful completion
//     - Event{Error: ...} if the provider rejected or failed mid-stream
//
// The channel is always closed by the provider after the terminating event.
// Callers must range until close or cancel via ctx.
package provider

import "context"

// Provider is the abstraction every LLM backend implements. Implementations
// must respect ctx cancellation and close the returned channel on completion.
type Provider interface {
	// Complete streams a completion for req. The returned channel receives zero or
	// more Delta events followed by exactly one terminating event (Usage on success,
	// Error on failure). Channel is closed after the terminating event or on context
	// cancellation.
	Complete(ctx context.Context, req CompleteRequest) (<-chan Event, error)

	// Name returns the provider identifier ("anthropic", "openai", "ollama").
	// Used for logs, metrics, and registry lookup.
	Name() string
}

// CompleteRequest is the normalised shape a handler builds before invoking a
// Provider. Model names are provider-specific strings — the router does not
// translate them.
type CompleteRequest struct {
	// Model is the provider-specific identifier, e.g. "claude-sonnet-4-6",
	// "gpt-4o", "llama3.2". Required.
	Model string

	// System is the system prompt, typically the rendered PromptContext.
	// May be empty for providers that don't support a separate system field
	// (implementations prepend it to the first user message).
	System string

	// Messages is the user/assistant conversation history. Required non-empty.
	Messages []Message

	// MaxTokens caps the response length. Zero means "use provider default".
	MaxTokens int
}

// Message is a single turn in the conversation. Role must be RoleUser or
// RoleAssistant.
type Message struct {
	Role    string
	Content string
}

// Conversation roles.
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Event is one element of the provider's output stream. Exactly one of Delta,
// Usage, or Error is populated per event in the terminal case; Delta may be ""
// on intermediate events that carry no content (rare, but observed with
// heartbeat frames).
type Event struct {
	// Delta is a partial text fragment. Present on intermediate events.
	Delta string

	// Usage is populated only on the terminating success event.
	Usage *UsageMetrics

	// Error is populated only on the terminating failure event.
	Error *ProviderError
}

// IsTerminal reports whether this event closes the stream (success or failure).
// Callers that need to branch on end-of-stream should prefer this over
// inspecting Usage/Error directly.
func (e Event) IsTerminal() bool {
	return e.Usage != nil || e.Error != nil
}

// UsageMetrics reports token consumption. InputTokens includes the system
// prompt. Providers that don't report usage leave these zero.
type UsageMetrics struct {
	InputTokens  int
	OutputTokens int
}

// ProviderError is the normalised error shape. Code is one of the Err* constants;
// Message is human-readable and SAFE TO LOG (does not include prompt content).
type ProviderError struct {
	Code    string
	Message string
}

// Error implements the error interface so ProviderError can be used wherever
// a Go error is expected.
func (e *ProviderError) Error() string {
	if e == nil {
		return ""
	}
	return e.Code + ": " + e.Message
}

// Canonical error codes. Providers translate native errors into these so the
// client sees a consistent shape regardless of backend.
const (
	// ErrRateLimit indicates the provider throttled the request (HTTP 429 or equivalent).
	ErrRateLimit = "rate_limit"
	// ErrInvalidKey indicates authentication failed (HTTP 401/403).
	ErrInvalidKey = "invalid_key"
	// ErrTimeout indicates the request exceeded a deadline.
	ErrTimeout = "timeout"
	// ErrProviderDown indicates a 5xx from the provider or network failure.
	ErrProviderDown = "provider_5xx"
	// ErrInvalidRequest indicates the request was malformed (HTTP 400).
	ErrInvalidRequest = "invalid_request"
	// ErrContextCanceled indicates the client disconnected or ctx was cancelled.
	ErrContextCanceled = "context_canceled"
	// ErrUnknown is the fallback when the native error can't be classified.
	ErrUnknown = "unknown"
)
