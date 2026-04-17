package provider

import "testing"

func TestEvent_IsTerminal(t *testing.T) {
	cases := []struct {
		name string
		e    Event
		want bool
	}{
		{"empty", Event{}, false},
		{"delta only", Event{Delta: "hi"}, false},
		{"usage", Event{Usage: &UsageMetrics{InputTokens: 10, OutputTokens: 5}}, true},
		{"error", Event{Error: &ProviderError{Code: ErrRateLimit, Message: "x"}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.e.IsTerminal(); got != tc.want {
				t.Errorf("IsTerminal() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestProviderError_Error(t *testing.T) {
	cases := []struct {
		name string
		e    *ProviderError
		want string
	}{
		{"nil", nil, ""},
		{"populated", &ProviderError{Code: ErrRateLimit, Message: "retry in 60s"}, "rate_limit: retry in 60s"},
		{"empty fields", &ProviderError{}, ": "},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.e.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestErrorCodes_Stable(t *testing.T) {
	// Guard against drift — error code strings are part of the public
	// contract with API consumers. Changing these requires an API version bump.
	expect := map[string]string{
		ErrRateLimit:       "rate_limit",
		ErrInvalidKey:      "invalid_key",
		ErrTimeout:         "timeout",
		ErrProviderDown:    "provider_5xx",
		ErrInvalidRequest:  "invalid_request",
		ErrContextCanceled: "context_canceled",
		ErrUnknown:         "unknown",
	}
	for code, want := range expect {
		if code != want {
			t.Errorf("error code changed: got %q want %q", code, want)
		}
	}
}
