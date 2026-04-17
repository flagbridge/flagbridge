package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AnthropicAPIVersion is the required header for the Messages API.
// Bumping this is an ADR-worthy decision — behaviour changes between versions.
const AnthropicAPIVersion = "2023-06-01"

// defaultAnthropicBaseURL is the production endpoint. Tests override via NewAnthropic.
const defaultAnthropicBaseURL = "https://api.anthropic.com"

// defaultMaxTokens is used when CompleteRequest.MaxTokens is 0.
// 4096 covers the longest expected Cmd+K response comfortably while capping cost.
const defaultMaxTokens = 4096

type anthropicProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewAnthropic returns a Provider backed by Anthropic's Messages API streaming
// endpoint. Pass empty baseURL to use the production default.
//
// The returned Provider uses a dedicated http.Client with no Timeout — streaming
// responses are inherently long-lived. Callers cancel via ctx.
func NewAnthropic(apiKey, baseURL string) Provider {
	if baseURL == "" {
		baseURL = defaultAnthropicBaseURL
	}
	return &anthropicProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{}, // Timeout: 0 — streaming
	}
}

// Name implements Provider.
func (p *anthropicProvider) Name() string { return "anthropic" }

// Complete implements Provider.Complete by posting to /v1/messages with stream=true
// and translating Anthropic's SSE frames into our canonical Event stream.
func (p *anthropicProvider) Complete(ctx context.Context, req CompleteRequest) (<-chan Event, error) {
	body := anthropicRequestBody{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		System:    req.System,
		Stream:    true,
	}
	if body.MaxTokens == 0 {
		body.MaxTokens = defaultMaxTokens
	}
	for _, m := range req.Messages {
		body.Messages = append(body.Messages, anthropicMessage{Role: m.Role, Content: m.Content})
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal anthropic request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build anthropic request: %w", err)
	}
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", AnthropicAPIVersion)
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("accept", "text/event-stream")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, classifyTransportError(err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, classifyHTTPError(resp.StatusCode, errBody)
	}

	ch := make(chan Event, 16)
	go p.streamResponse(ctx, resp, ch)
	return ch, nil
}

// streamResponse parses Anthropic's SSE frames and emits canonical Events.
//
// Anthropic emits one event per frame in a stream like:
//
//	event: message_start
//	data: {"type":"message_start","message":{"usage":{"input_tokens":25,"output_tokens":1}}}
//
//	event: content_block_delta
//	data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"Hello"}}
//
//	event: message_delta
//	data: {"type":"message_delta","delta":{...},"usage":{"output_tokens":15}}
//
//	event: message_stop
//	data: {"type":"message_stop"}
//
// We extract text from content_block_delta, accumulate input_tokens from
// message_start and output_tokens from message_delta, and emit the terminal
// Usage event on message_stop.
func (p *anthropicProvider) streamResponse(ctx context.Context, resp *http.Response, ch chan<- Event) {
	defer close(ch)
	defer resp.Body.Close()

	var inputTokens, outputTokens int
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // room for large JSON frames

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "" {
			continue
		}

		var frame anthropicFrame
		if err := json.Unmarshal([]byte(data), &frame); err != nil {
			// Anthropic emits `error` frames with different shape — try that.
			var errFrame anthropicErrorFrame
			if json.Unmarshal([]byte(data), &errFrame) == nil && errFrame.Type == "error" {
				emit(ctx, ch, Event{Error: &ProviderError{
					Code:    mapAnthropicErrorType(errFrame.Error.Type),
					Message: errFrame.Error.Message,
				}})
				return
			}
			continue // unknown frame, skip
		}

		switch frame.Type {
		case "message_start":
			if frame.Message != nil {
				inputTokens = frame.Message.Usage.InputTokens
				outputTokens = frame.Message.Usage.OutputTokens
			}
		case "content_block_delta":
			if frame.Delta != nil && frame.Delta.Type == "text_delta" && frame.Delta.Text != "" {
				if !emit(ctx, ch, Event{Delta: frame.Delta.Text}) {
					return
				}
			}
		case "message_delta":
			if frame.Usage != nil {
				outputTokens = frame.Usage.OutputTokens
			}
		case "message_stop":
			emit(ctx, ch, Event{Usage: &UsageMetrics{InputTokens: inputTokens, OutputTokens: outputTokens}})
			return
		case "error":
			if frame.ErrorPayload != nil {
				emit(ctx, ch, Event{Error: &ProviderError{
					Code:    mapAnthropicErrorType(frame.ErrorPayload.Type),
					Message: frame.ErrorPayload.Message,
				}})
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		emit(ctx, ch, Event{Error: &ProviderError{Code: ErrProviderDown, Message: "stream read: " + err.Error()}})
	}
}

// emit sends ev on ch unless ctx is done. Returns false if ctx was done.
func emit(ctx context.Context, ch chan<- Event, ev Event) bool {
	select {
	case <-ctx.Done():
		return false
	case ch <- ev:
		return true
	}
}

// classifyTransportError maps net/http.Client.Do errors into ProviderError.
func classifyTransportError(err error) error {
	if errors.Is(err, context.Canceled) {
		return &ProviderError{Code: ErrContextCanceled, Message: err.Error()}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &ProviderError{Code: ErrTimeout, Message: err.Error()}
	}
	return &ProviderError{Code: ErrProviderDown, Message: err.Error()}
}

// classifyHTTPError maps non-2xx responses into ProviderError. The upstream body
// may contain useful context (rate limit windows, invalid model names); we include
// up to 512 bytes of it in Message.
func classifyHTTPError(status int, body []byte) error {
	msg := strings.TrimSpace(string(body))
	if len(msg) > 512 {
		msg = msg[:512] + "..."
	}
	switch {
	case status == http.StatusTooManyRequests:
		return &ProviderError{Code: ErrRateLimit, Message: msg}
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return &ProviderError{Code: ErrInvalidKey, Message: msg}
	case status == http.StatusBadRequest:
		return &ProviderError{Code: ErrInvalidRequest, Message: msg}
	case status >= 500:
		return &ProviderError{Code: ErrProviderDown, Message: msg}
	default:
		return &ProviderError{Code: ErrUnknown, Message: fmt.Sprintf("status %d: %s", status, msg)}
	}
}

// mapAnthropicErrorType translates Anthropic's error.type values into canonical codes.
func mapAnthropicErrorType(t string) string {
	switch t {
	case "rate_limit_error", "overloaded_error":
		return ErrRateLimit
	case "authentication_error", "permission_error":
		return ErrInvalidKey
	case "invalid_request_error":
		return ErrInvalidRequest
	case "api_error":
		return ErrProviderDown
	default:
		return ErrUnknown
	}
}

// --- Wire types ---

type anthropicRequestBody struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Stream    bool               `json:"stream"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicFrame is the union of all SSE frame payloads we care about.
// Fields not relevant to a given frame type are nil.
type anthropicFrame struct {
	Type         string                 `json:"type"`
	Message      *anthropicMessageInfo  `json:"message,omitempty"`
	Delta        *anthropicDelta        `json:"delta,omitempty"`
	Usage        *anthropicUsage        `json:"usage,omitempty"`
	ErrorPayload *anthropicErrorPayload `json:"error,omitempty"`
}

type anthropicMessageInfo struct {
	Usage anthropicUsage `json:"usage"`
}

type anthropicDelta struct {
	Type string `json:"type"` // "text_delta" or "input_json_delta" etc.
	Text string `json:"text"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type anthropicErrorFrame struct {
	Type  string                `json:"type"`
	Error anthropicErrorPayload `json:"error"`
}

type anthropicErrorPayload struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
