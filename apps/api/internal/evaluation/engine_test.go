package evaluation

import (
	"encoding/json"
	"testing"
)

func raw(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func rawPtr(v any) *json.RawMessage {
	r := raw(v)
	return &r
}

// TestEvaluate_TableDriven is the canonical table-driven test covering the full
// resolution order: disabled → session_override → targeting_rule → rollout → default_value.
func TestEvaluate_TableDriven(t *testing.T) {
	// Bucket reference for "rollout-flag":
	//   user-g → bucket 0  (always in)
	//   user-b → bucket 12
	//   user-a → bucket 23
	//   user-h → bucket 38
	//   user-e → bucket 53 (out if rollout ≤ 53)
	//   user-d → bucket 93 (out if rollout ≤ 93)

	tests := []struct {
		name        string
		flagKey     string
		data        FlagData
		ctx         EvalContext
		wantReason  string
		wantValue   string
		wantRuleID  string
	}{
		// --- disabled ---
		{
			name:    "disabled flag returns default",
			flagKey: "my-flag",
			data:    FlagData{Enabled: false, DefaultValue: raw(false)},
			ctx:     EvalContext{},
			wantReason: ReasonDisabled,
			wantValue:  "false",
		},
		{
			name:    "disabled flag ignores rules",
			flagKey: "my-flag",
			data: FlagData{
				Enabled:      false,
				DefaultValue: raw("off"),
				Rules: []Rule{{ID: "r1", Priority: 1, Enabled: true, Value: raw("on")}},
			},
			ctx:        EvalContext{},
			wantReason: ReasonDisabled,
			wantValue:  `"off"`,
		},

		// --- session override ---
		{
			name:    "session override wins over everything",
			flagKey: "my-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules: []Rule{{
					ID: "r1", Priority: 1, Enabled: true,
					Conditions: []Condition{{Attribute: "country", Operator: OpEquals, Value: "BR"}},
					Value:      raw(true),
				}},
			},
			ctx: EvalContext{
				Attributes:      map[string]string{"country": "BR"},
				SessionOverride: rawPtr("forced-value"),
			},
			wantReason: ReasonSessionOverride,
			wantValue:  `"forced-value"`,
		},
		{
			name:    "session override boolean false",
			flagKey: "kill-switch",
			data:    FlagData{Enabled: true, DefaultValue: raw(true)},
			ctx:     EvalContext{SessionOverride: rawPtr(false)},
			wantReason: ReasonSessionOverride,
			wantValue:  "false",
		},
		{
			name:    "nil session override is ignored",
			flagKey: "my-flag",
			data:    FlagData{Enabled: true, DefaultValue: raw("default")},
			ctx:     EvalContext{SessionOverride: nil},
			wantReason: ReasonDefaultValue,
			wantValue:  `"default"`,
		},

		// --- targeting rules ---
		{
			name:    "single rule matches",
			flagKey: "waitlist",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules: []Rule{{
					ID: "r1", Priority: 1, Enabled: true,
					Conditions: []Condition{{Attribute: "country", Operator: OpEquals, Value: "BR"}},
					Value:      raw(true),
				}},
			},
			ctx:        EvalContext{Attributes: map[string]string{"country": "BR"}},
			wantReason: ReasonTargetingRule,
			wantValue:  "true",
			wantRuleID: "r1",
		},
		{
			name:    "rule does not match, falls to default",
			flagKey: "waitlist",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules: []Rule{{
					ID: "r1", Priority: 1, Enabled: true,
					Conditions: []Condition{{Attribute: "country", Operator: OpEquals, Value: "JP"}},
					Value:      raw(true),
				}},
			},
			ctx:        EvalContext{Attributes: map[string]string{"country": "BR"}},
			wantReason: ReasonDefaultValue,
			wantValue:  "false",
		},
		{
			name:    "higher-priority rule wins",
			flagKey: "feature-x",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw("default"),
				Rules: []Rule{
					{ID: "low", Priority: 10, Enabled: true,
						Conditions: []Condition{{Attribute: "locale", Operator: OpStartsWith, Value: "pt"}},
						Value:      raw("pt-value")},
					{ID: "high", Priority: 1, Enabled: true,
						Conditions: []Condition{{Attribute: "country", Operator: OpEquals, Value: "BR"}},
						Value:      raw("br-value")},
				},
			},
			ctx: EvalContext{Attributes: map[string]string{"country": "BR", "locale": "pt-BR"}},
			wantReason: ReasonTargetingRule,
			wantValue:  `"br-value"`,
			wantRuleID: "high",
		},
		{
			name:    "disabled rule is skipped",
			flagKey: "feature-x",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw("fallback"),
				Rules: []Rule{{
					ID: "disabled-rule", Priority: 1, Enabled: false,
					Conditions: []Condition{{Attribute: "country", Operator: OpEquals, Value: "BR"}},
					Value:      raw("should-not-match"),
				}},
			},
			ctx:        EvalContext{Attributes: map[string]string{"country": "BR"}},
			wantReason: ReasonDefaultValue,
			wantValue:  `"fallback"`,
		},
		{
			name:    "empty conditions rule always matches",
			flagKey: "catch-all",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules:        []Rule{{ID: "catch", Priority: 1, Enabled: true, Conditions: []Condition{}, Value: raw(true)}},
			},
			ctx:        EvalContext{},
			wantReason: ReasonTargetingRule,
			wantValue:  "true",
			wantRuleID: "catch",
		},
		{
			name:    "AND conditions: all must match",
			flagKey: "complex",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules: []Rule{{
					ID: "r1", Priority: 1, Enabled: true,
					Conditions: []Condition{
						{Attribute: "country", Operator: OpIn, Value: `["BR","MX"]`},
						{Attribute: "age", Operator: OpGTE, Value: "18"},
					},
					Value: raw(true),
				}},
			},
			ctx: EvalContext{Attributes: map[string]string{"country": "BR", "age": "20"}},
			wantReason: ReasonTargetingRule,
			wantValue:  "true",
			wantRuleID: "r1",
		},
		{
			name:    "AND conditions: one fails → default",
			flagKey: "complex",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules: []Rule{{
					ID: "r1", Priority: 1, Enabled: true,
					Conditions: []Condition{
						{Attribute: "country", Operator: OpIn, Value: `["BR","MX"]`},
						{Attribute: "age", Operator: OpGTE, Value: "18"},
					},
					Value: raw(true),
				}},
			},
			ctx:        EvalContext{Attributes: map[string]string{"country": "BR", "age": "16"}},
			wantReason: ReasonDefaultValue,
			wantValue:  "false",
		},
		{
			name:    "user_id targeting",
			flagKey: "beta",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rules: []Rule{{
					ID: "beta-users", Priority: 1, Enabled: true,
					Conditions: []Condition{{Attribute: "user_id", Operator: OpIn, Value: `["user-1","user-2"]`}},
					Value:      raw(true),
				}},
			},
			ctx:        EvalContext{UserID: "user-2"},
			wantReason: ReasonTargetingRule,
			wantValue:  "true",
			wantRuleID: "beta-users",
		},

		// --- rollout ---
		{
			name:    "rollout 100 always includes user",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				StateValue:   raw(true),
				Rollout:      100,
			},
			ctx:        EvalContext{UserID: "user-d"}, // bucket=93 < 100
			wantReason: ReasonRollout,
			wantValue:  "true",
		},
		{
			name:    "rollout 50 includes user with bucket 12",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				StateValue:   raw(true),
				Rollout:      50,
			},
			ctx:        EvalContext{UserID: "user-b"}, // bucket=12 < 50
			wantReason: ReasonRollout,
			wantValue:  "true",
		},
		{
			name:    "rollout 50 excludes user with bucket 53",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				Rollout:      50,
				// no StateValue — excluded user gets DefaultValue
			},
			ctx:        EvalContext{UserID: "user-e"}, // bucket=53 >= 50
			wantReason: ReasonDefaultValue,
			wantValue:  "false",
		},
		{
			name:    "rollout 0 disabled — no user in rollout",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				StateValue:   raw(true),
				Rollout:      0,
			},
			ctx:        EvalContext{UserID: "user-g"}, // bucket=0, but rollout=0 is disabled
			wantReason: ReasonDefaultValue,
			wantValue:  "true", // StateValue takes over
		},
		{
			name:    "rollout ignored when user_id is empty",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				StateValue:   raw(true),
				Rollout:      100,
			},
			ctx:        EvalContext{}, // no user_id
			wantReason: ReasonDefaultValue,
			wantValue:  "true", // StateValue
		},
		{
			name:    "rollout uses state_value when set",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw("off"),
				StateValue:   raw("on"),
				Rollout:      50,
			},
			ctx:        EvalContext{UserID: "user-b"}, // bucket=12 < 50
			wantReason: ReasonRollout,
			wantValue:  `"on"`,
		},
		{
			name:    "rollout falls back to default_value when state_value is nil",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw("base"),
				Rollout:      50,
			},
			ctx:        EvalContext{UserID: "user-b"}, // bucket=12 < 50
			wantReason: ReasonRollout,
			wantValue:  `"base"`,
		},
		{
			name:    "targeting rule beats rollout",
			flagKey: "rollout-flag",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw(false),
				StateValue:   raw(true),
				Rollout:      100,
				Rules: []Rule{{
					ID: "staff", Priority: 1, Enabled: true,
					Conditions: []Condition{{Attribute: "role", Operator: OpEquals, Value: "staff"}},
					Value:      raw("staff-value"),
				}},
			},
			ctx:        EvalContext{UserID: "user-d", Attributes: map[string]string{"role": "staff"}},
			wantReason: ReasonTargetingRule,
			wantValue:  `"staff-value"`,
			wantRuleID: "staff",
		},

		// --- state value fallback ---
		{
			name:    "no rules, state_value overrides default",
			flagKey: "hero-variant",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw("default"),
				StateValue:   raw("variant-b"),
			},
			ctx:        EvalContext{},
			wantReason: ReasonDefaultValue,
			wantValue:  `"variant-b"`,
		},
		{
			name:    "no rules, no state_value — returns default",
			flagKey: "hero-variant",
			data: FlagData{
				Enabled:      true,
				DefaultValue: raw("default"),
			},
			ctx:        EvalContext{},
			wantReason: ReasonDefaultValue,
			wantValue:  `"default"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Evaluate(tt.flagKey, tt.data, tt.ctx)
			if got.Reason != tt.wantReason {
				t.Errorf("reason: got %q, want %q", got.Reason, tt.wantReason)
			}
			if string(got.Value) != tt.wantValue {
				t.Errorf("value: got %s, want %s", got.Value, tt.wantValue)
			}
			if tt.wantRuleID != "" && got.RuleID != tt.wantRuleID {
				t.Errorf("rule_id: got %q, want %q", got.RuleID, tt.wantRuleID)
			}
			if got.FlagKey != tt.flagKey {
				t.Errorf("flag_key: got %q, want %q", got.FlagKey, tt.flagKey)
			}
		})
	}
}

// TestRolloutBucket verifies the hash function is deterministic and bounded.
func TestRolloutBucket(t *testing.T) {
	for i := 0; i < 100; i++ {
		b := rolloutBucket("user-abc", "some-flag")
		if b < 0 || b >= 100 {
			t.Fatalf("bucket out of range: %d", b)
		}
	}
	// Deterministic
	b1 := rolloutBucket("user-xyz", "flag-k")
	b2 := rolloutBucket("user-xyz", "flag-k")
	if b1 != b2 {
		t.Errorf("rolloutBucket is not deterministic: %d != %d", b1, b2)
	}
	// Different users get potentially different buckets
	bA := rolloutBucket("user-a", "rollout-flag") // known: 23
	bG := rolloutBucket("user-g", "rollout-flag") // known: 0
	if bA == bG {
		t.Errorf("expected different buckets for user-a and user-g, both got %d", bA)
	}
}

func TestEvaluate_FlagDisabled(t *testing.T) {
	result := Evaluate("my-flag", FlagData{
		Enabled:      false,
		DefaultValue: raw(false),
	}, EvalContext{})

	if result.Reason != ReasonDisabled {
		t.Errorf("expected reason %q, got %q", ReasonDisabled, result.Reason)
	}
	if string(result.Value) != "false" {
		t.Errorf("expected value false, got %s", result.Value)
	}
	if result.FlagKey != "my-flag" {
		t.Errorf("expected flag_key %q, got %q", "my-flag", result.FlagKey)
	}
}

func TestEvaluate_EnabledNoRules_DefaultValue(t *testing.T) {
	result := Evaluate("my-flag", FlagData{
		Enabled:      true,
		DefaultValue: raw(true),
		Rules:        nil,
	}, EvalContext{})

	if result.Reason != ReasonDefaultValue {
		t.Errorf("expected reason %q, got %q", ReasonDefaultValue, result.Reason)
	}
	if string(result.Value) != "true" {
		t.Errorf("expected value true, got %s", result.Value)
	}
}

func TestEvaluate_EnabledNoRules_StateValueOverridesDefault(t *testing.T) {
	result := Evaluate("hero-variant", FlagData{
		Enabled:      true,
		DefaultValue: raw("default"),
		StateValue:   raw("variant-b"),
		Rules:        nil,
	}, EvalContext{})

	if result.Reason != ReasonDefaultValue {
		t.Errorf("expected reason %q, got %q", ReasonDefaultValue, result.Reason)
	}
	if string(result.Value) != `"variant-b"` {
		t.Errorf("expected value %q, got %s", `"variant-b"`, result.Value)
	}
}

func TestEvaluate_OneRuleMatches(t *testing.T) {
	result := Evaluate("waitlist-open", FlagData{
		Enabled:      true,
		DefaultValue: raw(false),
		Rules: []Rule{
			{
				ID:       "rule-1",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "country", Operator: OpEquals, Value: "BR"},
				},
				Value: raw(true),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{"country": "BR"},
	})

	if result.Reason != ReasonTargetingRule {
		t.Errorf("expected reason %q, got %q", ReasonTargetingRule, result.Reason)
	}
	if result.RuleID != "rule-1" {
		t.Errorf("expected rule_id %q, got %q", "rule-1", result.RuleID)
	}
	if string(result.Value) != "true" {
		t.Errorf("expected value true, got %s", result.Value)
	}
}

func TestEvaluate_MultipleRules_FirstMatchByPriority(t *testing.T) {
	result := Evaluate("feature-x", FlagData{
		Enabled:      true,
		DefaultValue: raw("default"),
		Rules: []Rule{
			{
				ID:       "rule-low-priority",
				Priority: 10,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "locale", Operator: OpStartsWith, Value: "pt"},
				},
				Value: raw("pt-value"),
			},
			{
				ID:       "rule-high-priority",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "country", Operator: OpEquals, Value: "BR"},
				},
				Value: raw("br-value"),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{
			"country": "BR",
			"locale":  "pt-BR",
		},
	})

	if result.RuleID != "rule-high-priority" {
		t.Errorf("expected first match by priority, got rule_id %q", result.RuleID)
	}
	if string(result.Value) != `"br-value"` {
		t.Errorf("expected %q, got %s", `"br-value"`, result.Value)
	}
}

func TestEvaluate_NoRulesMatch_ReturnsDefault(t *testing.T) {
	result := Evaluate("feature-x", FlagData{
		Enabled:      true,
		DefaultValue: raw("fallback"),
		Rules: []Rule{
			{
				ID:       "rule-1",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "country", Operator: OpEquals, Value: "JP"},
				},
				Value: raw("jp-only"),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{"country": "BR"},
	})

	if result.Reason != ReasonDefaultValue {
		t.Errorf("expected reason %q, got %q", ReasonDefaultValue, result.Reason)
	}
}

func TestEvaluate_DisabledRule_Skipped(t *testing.T) {
	result := Evaluate("feature-x", FlagData{
		Enabled:      true,
		DefaultValue: raw("fallback"),
		Rules: []Rule{
			{
				ID:       "disabled-rule",
				Priority: 1,
				Enabled:  false,
				Conditions: []Condition{
					{Attribute: "country", Operator: OpEquals, Value: "BR"},
				},
				Value: raw("should-not-match"),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{"country": "BR"},
	})

	if result.Reason != ReasonDefaultValue {
		t.Errorf("disabled rule should be skipped, got reason %q", result.Reason)
	}
}

func TestEvaluate_ComplexANDConditions(t *testing.T) {
	result := Evaluate("complex-flag", FlagData{
		Enabled:      true,
		DefaultValue: raw(false),
		Rules: []Rule{
			{
				ID:       "complex-rule",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "country", Operator: OpIn, Value: `["BR","MX","AR"]`},
					{Attribute: "locale", Operator: OpStartsWith, Value: "pt"},
					{Attribute: "age", Operator: OpGTE, Value: "18"},
				},
				Value: raw(true),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{
			"country": "BR",
			"locale":  "pt-BR",
			"age":     "25",
		},
	})

	if result.Reason != ReasonTargetingRule {
		t.Errorf("all AND conditions should match, got reason %q", result.Reason)
	}
}

func TestEvaluate_ComplexANDConditions_OneFails(t *testing.T) {
	result := Evaluate("complex-flag", FlagData{
		Enabled:      true,
		DefaultValue: raw(false),
		Rules: []Rule{
			{
				ID:       "complex-rule",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "country", Operator: OpIn, Value: `["BR","MX","AR"]`},
					{Attribute: "locale", Operator: OpStartsWith, Value: "pt"},
					{Attribute: "age", Operator: OpGTE, Value: "18"},
				},
				Value: raw(true),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{
			"country": "BR",
			"locale":  "pt-BR",
			"age":     "16",
		},
	})

	if result.Reason != ReasonDefaultValue {
		t.Errorf("one condition fails, should return default, got reason %q", result.Reason)
	}
}

func TestEvaluate_EmptyConditions_RuleAlwaysMatches(t *testing.T) {
	result := Evaluate("catch-all", FlagData{
		Enabled:      true,
		DefaultValue: raw(false),
		Rules: []Rule{
			{
				ID:         "catch-all-rule",
				Priority:   1,
				Enabled:    true,
				Conditions: []Condition{},
				Value:      raw(true),
			},
		},
	}, EvalContext{})

	if result.Reason != ReasonTargetingRule {
		t.Errorf("empty conditions should always match, got reason %q", result.Reason)
	}
}

func TestEvaluate_UserIDTargeting(t *testing.T) {
	result := Evaluate("beta-feature", FlagData{
		Enabled:      true,
		DefaultValue: raw(false),
		Rules: []Rule{
			{
				ID:       "beta-users",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "user_id", Operator: OpIn, Value: `["user-1","user-2","user-3"]`},
				},
				Value: raw(true),
			},
		},
	}, EvalContext{
		UserID: "user-2",
	})

	if result.Reason != ReasonTargetingRule {
		t.Errorf("user_id should be targetable, got reason %q", result.Reason)
	}
}

func TestEvaluate_StringValue(t *testing.T) {
	result := Evaluate("hero-variant", FlagData{
		Enabled:      true,
		DefaultValue: raw("control"),
		Rules: []Rule{
			{
				ID:       "variant-a-rule",
				Priority: 1,
				Enabled:  true,
				Conditions: []Condition{
					{Attribute: "cohort", Operator: OpEquals, Value: "experiment"},
				},
				Value: raw("variant-a"),
			},
		},
	}, EvalContext{
		Attributes: map[string]string{"cohort": "experiment"},
	})

	if string(result.Value) != `"variant-a"` {
		t.Errorf("expected string value %q, got %s", `"variant-a"`, result.Value)
	}
}
