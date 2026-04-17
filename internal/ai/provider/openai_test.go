package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const openaiSSEScript = `data: {"id":"chatcmpl-1","choices":[{"index":0,"delta":{"role":"assistant","content":""}}]}

data: {"id":"chatcmpl-1","choices":[{"index":0,"delta":{"content":"Hello"}}]}

data: {"id":"chatcmpl-1","choices":[{"index":0,"delta":{"content":", "}}]}

data: {"id":"chatcmpl-1","choices":[{"index":0,"delta":{"content":"world!"}}]}

data: {"id":"chatcmpl-1","choices":[{"index":0,"finish_reason":"stop","delta":{}}]}

data: {"id":"chatcmpl-1","choices":[],"usage":{"prompt_tokens":25,"completion_tokens":15,"total_tokens":40}}

data: [DONE]

`

func TestOpenAI_HappyPath(t *testing.T) {
	var capturedBody openaiRequestBody
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("authorization"); got != "Bearer test-key" {
			t.Errorf("authorization = %q, want %q", got, "Bearer test-key")
		}
		_ = decodeJSONBody(r.Body, &capturedBody)
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(openaiSSEScript))
	}))
	defer srv.Close()

	p := NewOpenAI("test-key", srv.URL)
	ch, err := p.Complete(context.Background(), CompleteRequest{
		Model:    "gpt-4o",
		System:   "flag assistant",
		Messages: []Message{{Role: RoleUser, Content: "create a flag"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	var text strings.Builder
	var terminal Event
	for ev := range ch {
		if ev.IsTerminal() {
			terminal = ev
			continue
		}
		text.WriteString(ev.Delta)
	}

	if got := text.String(); got != "Hello, world!" {
		t.Errorf("assembled text = %q, want %q", got, "Hello, world!")
	}
	if terminal.Usage == nil {
		t.Fatal("expected Usage terminal")
	}
	if terminal.Usage.InputTokens != 25 || terminal.Usage.OutputTokens != 15 {
		t.Errorf("Usage = %+v, want {25, 15}", terminal.Usage)
	}

	// Verify System was prepended as first message.
	if len(capturedBody.Messages) != 2 {
		t.Fatalf("expected 2 messages (system + user), got %d", len(capturedBody.Messages))
	}
	if capturedBody.Messages[0].Role != "system" {
		t.Errorf("first message role = %q, want system", capturedBody.Messages[0].Role)
	}
	if capturedBody.StreamOptions == nil || !capturedBody.StreamOptions.IncludeUsage {
		t.Errorf("stream_options.include_usage must be true")
	}
}

func TestOpenAI_InvalidKeyViaErrorEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"Incorrect API key provided","type":"invalid_request_error","code":"invalid_api_key"}}`))
	}))
	defer srv.Close()

	p := NewOpenAI("bad-key", srv.URL)
	_, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	var pe *ProviderError
	if !errorsAs(err, &pe) || pe.Code != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestOpenAI_RateLimitViaCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"Rate limit reached","type":"rate_limit_error","code":"rate_limit_exceeded"}}`))
	}))
	defer srv.Close()

	p := NewOpenAI("k", srv.URL)
	_, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	var pe *ProviderError
	if !errorsAs(err, &pe) || pe.Code != ErrRateLimit {
		t.Errorf("expected ErrRateLimit, got %v", err)
	}
}

func TestOpenAI_ContextLengthExceeded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"Context too long","type":"invalid_request_error","code":"context_length_exceeded"}}`))
	}))
	defer srv.Close()

	p := NewOpenAI("k", srv.URL)
	_, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	var pe *ProviderError
	if !errorsAs(err, &pe) || pe.Code != ErrInvalidRequest {
		t.Errorf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestOpenAI_MidStreamError(t *testing.T) {
	script := `data: {"id":"1","choices":[{"index":0,"delta":{"content":"Hel"}}]}

data: {"error":{"message":"server error","type":"api_error","code":""}}

data: [DONE]

`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(script))
	}))
	defer srv.Close()

	p := NewOpenAI("k", srv.URL)
	ch, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	var deltas []string
	var terminal *ProviderError
	for ev := range ch {
		if ev.Error != nil {
			terminal = ev.Error
			continue
		}
		if ev.Delta != "" {
			deltas = append(deltas, ev.Delta)
		}
	}

	if len(deltas) != 1 || deltas[0] != "Hel" {
		t.Errorf("pre-error deltas = %v, want [Hel]", deltas)
	}
	if terminal == nil || terminal.Code != ErrProviderDown {
		t.Errorf("expected ErrProviderDown terminal, got %+v", terminal)
	}
}

func TestOpenAI_NoSystemMessage_OmitsSystemEntry(t *testing.T) {
	var captured openaiRequestBody
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = decodeJSONBody(r.Body, &captured)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	p := NewOpenAI("k", srv.URL)
	ch, err := p.Complete(context.Background(), CompleteRequest{
		Model:    "x",
		Messages: []Message{{Role: RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	for range ch {
	}

	if len(captured.Messages) != 1 {
		t.Errorf("without System, expected 1 message; got %d", len(captured.Messages))
	}
	if captured.Messages[0].Role != RoleUser {
		t.Errorf("sole message role = %q, want user", captured.Messages[0].Role)
	}
}
