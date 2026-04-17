package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// anthropicSSEScript emits a realistic Messages API streaming response.
const anthropicSSEScript = `event: message_start
data: {"type":"message_start","message":{"usage":{"input_tokens":25,"output_tokens":1}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":", "}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"world!"}}

event: content_block_stop
data: {"type":"content_block_stop","index":0}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":15}}

event: message_stop
data: {"type":"message_stop"}

`

func anthropicMockServer(t *testing.T, script string, wantHeaders map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range wantHeaders {
			if got := r.Header.Get(k); got != v {
				t.Errorf("header %s = %q, want %q", k, got, v)
			}
		}
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(script))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}))
}

func TestAnthropic_HappyPath(t *testing.T) {
	srv := anthropicMockServer(t, anthropicSSEScript, map[string]string{
		"x-api-key":         "test-key",
		"anthropic-version": AnthropicAPIVersion,
		"accept":            "text/event-stream",
	})
	defer srv.Close()

	p := NewAnthropic("test-key", srv.URL)

	ch, err := p.Complete(context.Background(), CompleteRequest{
		Model:    "claude-sonnet-4-6",
		System:   "You are a feature flag assistant.",
		Messages: []Message{{Role: RoleUser, Content: "create a dark-mode flag"}},
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
	if terminal.Error != nil {
		t.Fatalf("got error terminal: %+v", terminal.Error)
	}
	if terminal.Usage == nil {
		t.Fatal("expected Usage terminal")
	}
	if terminal.Usage.InputTokens != 25 {
		t.Errorf("InputTokens = %d, want 25", terminal.Usage.InputTokens)
	}
	if terminal.Usage.OutputTokens != 15 {
		t.Errorf("OutputTokens = %d, want 15", terminal.Usage.OutputTokens)
	}
}

func TestAnthropic_RateLimitHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"rate_limit_error","message":"retry in 60s"}}`))
	}))
	defer srv.Close()

	p := NewAnthropic("k", srv.URL)
	_, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	var pe *ProviderError
	if !errorsAs(err, &pe) {
		t.Fatalf("expected *ProviderError, got %T: %v", err, err)
	}
	if pe.Code != ErrRateLimit {
		t.Errorf("Code = %q, want %q", pe.Code, ErrRateLimit)
	}
}

func TestAnthropic_InvalidKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`unauthorized`))
	}))
	defer srv.Close()

	p := NewAnthropic("bad", srv.URL)
	_, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	var pe *ProviderError
	if !errorsAs(err, &pe) || pe.Code != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestAnthropic_ProviderDown5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	p := NewAnthropic("k", srv.URL)
	_, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	var pe *ProviderError
	if !errorsAs(err, &pe) || pe.Code != ErrProviderDown {
		t.Errorf("expected ErrProviderDown, got %v", err)
	}
}

func TestAnthropic_MidStreamError(t *testing.T) {
	script := `event: message_start
data: {"type":"message_start","message":{"usage":{"input_tokens":10,"output_tokens":0}}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hel"}}

event: error
data: {"type":"error","error":{"type":"overloaded_error","message":"server overloaded"}}

`
	srv := anthropicMockServer(t, script, nil)
	defer srv.Close()

	p := NewAnthropic("k", srv.URL)
	ch, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	var deltas []string
	var terminalErr *ProviderError
	for ev := range ch {
		if ev.Error != nil {
			terminalErr = ev.Error
			continue
		}
		if ev.Delta != "" {
			deltas = append(deltas, ev.Delta)
		}
	}

	if len(deltas) != 1 || deltas[0] != "Hel" {
		t.Errorf("pre-error deltas = %v, want [Hel]", deltas)
	}
	if terminalErr == nil {
		t.Fatal("expected terminal error")
	}
	if terminalErr.Code != ErrRateLimit {
		t.Errorf("Code = %q, want %q (overloaded maps to rate_limit)", terminalErr.Code, ErrRateLimit)
	}
}

func TestAnthropic_ContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		// Emit one chunk then hang.
		_, _ = w.Write([]byte("event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\"a\"}}\n\n"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		// Hang: select forever so the test can cancel.
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	p := NewAnthropic("k", srv.URL)
	ch, err := p.Complete(ctx, CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	<-ch // receive first delta
	cancel()

	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("channel did not close after cancel — streamResponse leaked")
	}
}

func TestAnthropic_DefaultsMaxTokens(t *testing.T) {
	var receivedMaxTokens int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body anthropicRequestBody
		_ = decodeJSON(r, &body)
		receivedMaxTokens = body.MaxTokens
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"))
	}))
	defer srv.Close()

	p := NewAnthropic("k", srv.URL)
	ch, err := p.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	for range ch {
	}

	if receivedMaxTokens != defaultMaxTokens {
		t.Errorf("MaxTokens = %d, want default %d", receivedMaxTokens, defaultMaxTokens)
	}
}

// --- helpers shared across provider tests ---

// errorsAs wraps errors.As to keep call sites short.
func errorsAs(err error, target **ProviderError) bool {
	if err == nil {
		return false
	}
	if pe, ok := err.(*ProviderError); ok {
		*target = pe
		return true
	}
	return false
}

func decodeJSON(r *http.Request, v any) error {
	return decodeJSONBody(r.Body, v)
}
