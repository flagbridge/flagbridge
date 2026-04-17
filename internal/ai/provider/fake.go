package provider

import (
	"context"
	"time"
)

// FakeProvider is a scripted in-memory Provider used in tests of upstream
// components (handlers, middleware, rate limiter). Tests configure the script
// to emit specific deltas, usage, or errors, then assert on caller behaviour.
//
// FakeProvider is NOT thread-safe for configuration — set fields before calling
// Complete. Multiple concurrent Complete calls on the same FakeProvider are
// safe (each gets its own goroutine and channel).
type FakeProvider struct {
	// NameValue is returned by Name(). Defaults to "fake".
	NameValue string

	// Deltas are emitted one-by-one as Event{Delta: ...}.
	Deltas []string

	// Usage is the terminal usage event. If nil and ErrorAt == -1, no terminal
	// event is emitted (caller relies on channel close alone).
	Usage *UsageMetrics

	// ErrorAt is the index in Deltas at which to emit ErrorEvent and stop.
	// -1 means "never". 0 means "before any delta". len(Deltas) means "after all deltas".
	ErrorAt int

	// ErrorEvent is the error to emit when ErrorAt is reached.
	ErrorEvent *ProviderError

	// DelayPerChunk simulates network latency between deltas.
	DelayPerChunk time.Duration
}

// NewFake constructs a FakeProvider with ErrorAt defaulted to -1 (no error).
func NewFake(deltas ...string) *FakeProvider {
	return &FakeProvider{
		NameValue: "fake",
		Deltas:    deltas,
		ErrorAt:   -1,
	}
}

// Name implements Provider.
func (f *FakeProvider) Name() string {
	if f.NameValue == "" {
		return "fake"
	}
	return f.NameValue
}

// Complete implements Provider. Emits the scripted deltas, then either an error
// or a usage event, then closes the channel.
func (f *FakeProvider) Complete(ctx context.Context, req CompleteRequest) (<-chan Event, error) {
	ch := make(chan Event, len(f.Deltas)+1)

	go func() {
		defer close(ch)

		for i, delta := range f.Deltas {
			if f.ErrorAt == i && f.ErrorEvent != nil {
				select {
				case <-ctx.Done():
					return
				case ch <- Event{Error: f.ErrorEvent}:
				}
				return
			}

			if f.DelayPerChunk > 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(f.DelayPerChunk):
				}
			}

			select {
			case <-ctx.Done():
				return
			case ch <- Event{Delta: delta}:
			}
		}

		if f.ErrorAt == len(f.Deltas) && f.ErrorEvent != nil {
			select {
			case <-ctx.Done():
				return
			case ch <- Event{Error: f.ErrorEvent}:
			}
			return
		}

		if f.Usage != nil {
			select {
			case <-ctx.Done():
				return
			case ch <- Event{Usage: f.Usage}:
			}
		}
	}()

	return ch, nil
}
