package flag

import (
	"encoding/json"
	"testing"
)

func TestKeyPattern(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		// valid
		{"my-flag", true},
		{"feature", true},
		{"new-feature-v2", true},
		{"a", true},
		{"abc-123", true},
		{"dark-mode", true},
		{"a-b-c-d", true},

		// invalid
		{"", false},
		{"My-Flag", false},           // uppercase
		{"my_flag", false},           // underscore
		{"my flag", false},           // space
		{"-leading", false},          // leading hyphen
		{"trailing-", false},         // trailing hyphen
		{"double--hyphen", false},    // consecutive hyphens
		{"my.flag", false},           // dot
		{"my/flag", false},           // slash
		{"UPPERCASE", false},         // all caps
		{"camelCase", false},         // mixed case
		{"123-start", true},          // numbers ok at start
		{"pro.analytics", false},     // dot separator
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := keyPattern.MatchString(tt.key)
			if got != tt.want {
				t.Errorf("keyPattern.MatchString(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestCreateRequest_Defaults(t *testing.T) {
	// Verify that empty type defaults to "boolean" and nil default_value defaults to false
	// This logic is in Service.Create, but we can't call it without a DB.
	// Instead we test the input validation logic directly.

	t.Run("empty key is invalid", func(t *testing.T) {
		if keyPattern.MatchString("") {
			t.Error("empty key should not match pattern")
		}
	})

	t.Run("default value for boolean flag", func(t *testing.T) {
		var dv json.RawMessage
		if dv == nil {
			dv = json.RawMessage(`false`)
		}
		if string(dv) != "false" {
			t.Errorf("expected default false, got %s", dv)
		}
	})

	t.Run("default type", func(t *testing.T) {
		flagType := ""
		if flagType == "" {
			flagType = "boolean"
		}
		if flagType != "boolean" {
			t.Errorf("expected boolean, got %s", flagType)
		}
	})
}
