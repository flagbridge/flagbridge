package provider

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestFakeProvider_ScriptedCompletion(t *testing.T) {
	f := NewFake("Hello", ", ", "world", "!")
	f.Usage = &UsageMetrics{InputTokens: 10, OutputTokens: 4}

	ch, err := f.Complete(context.Background(), CompleteRequest{
		Model:    "fake-model",
		Messages: []Message{{Role: RoleUser, Content: "hi"}},
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
		t.Fatal("expected terminal Usage event")
	}
	if terminal.Usage.OutputTokens != 4 {
		t.Errorf("OutputTokens = %d, want 4", terminal.Usage.OutputTokens)
	}
}

func TestFakeProvider_ErrorAtStart(t *testing.T) {
	f := NewFake("never", "emitted")
	f.ErrorAt = 0
	f.ErrorEvent = &ProviderError{Code: ErrInvalidKey, Message: "bad api key"}

	ch, err := f.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	var deltas []string
	var gotErr *ProviderError
	for ev := range ch {
		if ev.Error != nil {
			gotErr = ev.Error
		} else if ev.Delta != "" {
			deltas = append(deltas, ev.Delta)
		}
	}

	if len(deltas) != 0 {
		t.Errorf("expected no deltas before error, got %v", deltas)
	}
	if gotErr == nil || gotErr.Code != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey terminal, got %+v", gotErr)
	}
}

func TestFakeProvider_ErrorMidStream(t *testing.T) {
	f := NewFake("first", "second", "never-emitted")
	f.ErrorAt = 2 // before the third delta
	f.ErrorEvent = &ProviderError{Code: ErrProviderDown, Message: "server died"}

	ch, _ := f.Complete(context.Background(), CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})

	var deltas []string
	for ev := range ch {
		if ev.Delta != "" {
			deltas = append(deltas, ev.Delta)
		}
	}

	if len(deltas) != 2 || deltas[0] != "first" || deltas[1] != "second" {
		t.Errorf("expected [first second], got %v", deltas)
	}
}

func TestFakeProvider_ContextCancel(t *testing.T) {
	f := NewFake("a", "b", "c", "d")
	f.DelayPerChunk = 50 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	ch, _ := f.Complete(ctx, CompleteRequest{Model: "x", Messages: []Message{{Role: RoleUser, Content: "hi"}}})

	// Read one event then cancel.
	<-ch
	cancel()

	// Drain — channel should close quickly without emitting remaining deltas.
	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel did not close after cancel — producer leaked")
	}
}

func TestFakeProvider_Name(t *testing.T) {
	if got := NewFake().Name(); got != "fake" {
		t.Errorf("default name = %q, want %q", got, "fake")
	}
	f := &FakeProvider{NameValue: "custom"}
	if got := f.Name(); got != "custom" {
		t.Errorf("custom name = %q, want custom", got)
	}
}
