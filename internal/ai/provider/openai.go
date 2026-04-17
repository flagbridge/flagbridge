package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// defaultOpenAIBaseURL is the production endpoint. Tests override via NewOpenAI.
// Ollama (Sprint 7) is wire-compatible and plugs in by pointing this at the local
// Ollama server — see ADR-0006.
const defaultOpenAIBaseURL = "https://api.openai.com"

type openaiProvider struct {
	apiKey  string
	baseURL string
	name    string
	client  *http.Client
}

// NewOpenAI returns a Provider backed by OpenAI's Chat Completions streaming
// endpoint. Pass empty baseURL to use the production default.
func NewOpenAI(apiKey, baseURL string) Provider {
	if baseURL == "" {
		baseURL = defaultOpenAIBaseURL
	}
	return &openaiProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		name:    "openai",
		client:  &http.Client{}, // Timeout: 0 — streaming
	}
}

// Name implements Provider.
func (p *openaiProvider) Name() string { return p.name }

// Complete implements Provider.Complete by posting to /v1/chat/completions with
// stream=true and stream_options.include_usage=true so we get a terminal usage
// event before the [DONE] sentinel.
func (p *openaiProvider) Complete(ctx context.Context, req CompleteRequest) (<-chan Event, error) {
	messages := make([]openaiMessage, 0, len(req.Messages)+1)
	if req.System != "" {
		messages = append(messages, openaiMessage{Role: "system", Content: req.System})
	}
	for _, m := range req.Messages {
		messages = append(messages, openaiMessage{Role: m.Role, Content: m.Content})
	}

	body := openaiRequestBody{
		Model:         req.Model,
		Stream:        true,
		StreamOptions: &openaiStreamOptions{IncludeUsage: true},
		Messages:      messages,
	}
	if req.MaxTokens > 0 {
		body.MaxTokens = req.MaxTokens
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal openai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build openai request: %w", err)
	}
	httpReq.Header.Set("authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("accept", "text/event-stream")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, classifyTransportError(err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, classifyOpenAIHTTPError(resp.StatusCode, errBody)
	}

	ch := make(chan Event, 16)
	go p.streamResponse(ctx, resp, ch)
	return ch, nil
}

// streamResponse parses OpenAI SSE frames. Each frame is `data: <json>` or the
// terminal `data: [DONE]`. Intermediate frames carry choices[0].delta.content;
// the usage frame (emitted because stream_options.include_usage=true) carries
// prompt_tokens + completion_tokens.
func (p *openaiProvider) streamResponse(ctx context.Context, resp *http.Response, ch chan<- Event) {
	defer close(ch)
	defer resp.Body.Close()

	var usage *UsageMetrics
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

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

		if data == "[DONE]" {
			if usage != nil {
				emit(ctx, ch, Event{Usage: usage})
			}
			return
		}

		// Try error frame first (different shape).
		var errFrame openaiErrorFrame
		if json.Unmarshal([]byte(data), &errFrame) == nil && errFrame.Error != nil {
			emit(ctx, ch, Event{Error: &ProviderError{
				Code:    mapOpenAIErrorType(errFrame.Error.Type, errFrame.Error.Code),
				Message: errFrame.Error.Message,
			}})
			return
		}

		var frame openaiFrame
		if err := json.Unmarshal([]byte(data), &frame); err != nil {
			continue // unknown shape, skip
		}

		if frame.Usage != nil {
			usage = &UsageMetrics{
				InputTokens:  frame.Usage.PromptTokens,
				OutputTokens: frame.Usage.CompletionTokens,
			}
		}

		for _, c := range frame.Choices {
			if c.Delta.Content != "" {
				if !emit(ctx, ch, Event{Delta: c.Delta.Content}) {
					return
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		emit(ctx, ch, Event{Error: &ProviderError{Code: ErrProviderDown, Message: "stream read: " + err.Error()}})
	}
}

// classifyOpenAIHTTPError first tries to parse the OpenAI error envelope
// ({"error":{"type":"...","message":"..."}}) then falls back to the generic
// status-code classifier.
func classifyOpenAIHTTPError(status int, body []byte) error {
	var env openaiErrorFrame
	if json.Unmarshal(body, &env) == nil && env.Error != nil {
		return &ProviderError{
			Code:    mapOpenAIErrorType(env.Error.Type, env.Error.Code),
			Message: env.Error.Message,
		}
	}
	return classifyHTTPError(status, body)
}

// mapOpenAIErrorType translates OpenAI's error type/code pair into canonical codes.
// OpenAI uses "type" for category (e.g. "invalid_request_error") and "code" for
// specifics (e.g. "invalid_api_key"). Code is usually more specific when present.
func mapOpenAIErrorType(errType, errCode string) string {
	switch errCode {
	case "invalid_api_key", "account_deactivated":
		return ErrInvalidKey
	case "rate_limit_exceeded", "insufficient_quota":
		return ErrRateLimit
	case "context_length_exceeded", "model_not_found":
		return ErrInvalidRequest
	}
	switch errType {
	case "invalid_request_error":
		return ErrInvalidRequest
	case "authentication_error":
		return ErrInvalidKey
	case "rate_limit_error":
		return ErrRateLimit
	case "server_error", "api_error":
		return ErrProviderDown
	}
	return ErrUnknown
}

// --- Wire types ---

type openaiRequestBody struct {
	Model         string               `json:"model"`
	Messages      []openaiMessage      `json:"messages"`
	Stream        bool                 `json:"stream"`
	StreamOptions *openaiStreamOptions `json:"stream_options,omitempty"`
	MaxTokens     int                  `json:"max_tokens,omitempty"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiStreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type openaiFrame struct {
	Choices []openaiChoice `json:"choices"`
	Usage   *openaiUsage   `json:"usage,omitempty"`
}

type openaiChoice struct {
	Index        int         `json:"index"`
	Delta        openaiDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type openaiDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type openaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openaiErrorFrame struct {
	Error *openaiErrorPayload `json:"error,omitempty"`
}

type openaiErrorPayload struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
