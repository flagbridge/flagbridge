package evaluation

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"sort"
)

// Condition represents a single targeting condition.
type Condition struct {
	Attribute string   `json:"attribute"`
	Operator  Operator `json:"operator"`
	Value     string   `json:"value"`
}

// Rule represents a targeting rule with conditions and a result value.
type Rule struct {
	ID         string          `json:"id"`
	Priority   int             `json:"priority"`
	Conditions []Condition     `json:"conditions"`
	Value      json.RawMessage `json:"value"`
	Enabled    bool            `json:"enabled"`
}

// FlagData contains all data needed to evaluate a flag.
type FlagData struct {
	Enabled      bool
	DefaultValue json.RawMessage
	StateValue   json.RawMessage // value from flag_state (overrides default when no rule matches)
	Rules        []Rule
	Rollout      int // 0-100 percentage; 0 = disabled (no rollout)
}

// EvalContext is the context provided by the SDK for evaluation.
type EvalContext struct {
	UserID          string            `json:"user_id"`
	Attributes      map[string]string `json:"attributes"`
	SessionOverride *json.RawMessage  `json:"session_override,omitempty"`
}

// Result is the output of a flag evaluation.
type Result struct {
	FlagKey string          `json:"flag_key"`
	Value   json.RawMessage `json:"value"`
	Reason  string          `json:"reason"`
	RuleID  string          `json:"rule_id,omitempty"`
}

const (
	ReasonDisabled        = "disabled"
	ReasonSessionOverride = "session_override"
	ReasonTargetingRule   = "targeting_rule"
	ReasonRollout         = "rollout"
	ReasonDefaultValue    = "default_value"
)

// Evaluate resolves a flag value given the flag data and evaluation context.
// Resolution order: disabled → session_override → targeting_rule → rollout → default_value
func Evaluate(flagKey string, data FlagData, ctx EvalContext) Result {
	if !data.Enabled {
		return Result{
			FlagKey: flagKey,
			Value:   data.DefaultValue,
			Reason:  ReasonDisabled,
		}
	}

	// Session override takes highest precedence (testing API)
	if ctx.SessionOverride != nil && len(*ctx.SessionOverride) > 0 {
		return Result{
			FlagKey: flagKey,
			Value:   *ctx.SessionOverride,
			Reason:  ReasonSessionOverride,
		}
	}

	// Sort rules by priority ASC (lower number = higher priority)
	rules := make([]Rule, len(data.Rules))
	copy(rules, data.Rules)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if matchesAllConditions(rule.Conditions, ctx) {
			return Result{
				FlagKey: flagKey,
				Value:   rule.Value,
				Reason:  ReasonTargetingRule,
				RuleID:  rule.ID,
			}
		}
	}

	// Rollout percentage: hash user_id+flag_key → bucket 0-99
	// User is in rollout if bucket < Rollout (e.g. Rollout=10 → ~10% of users)
	if data.Rollout > 0 && ctx.UserID != "" {
		bucket := rolloutBucket(ctx.UserID, flagKey)
		if bucket < data.Rollout {
			value := data.DefaultValue
			if data.StateValue != nil && len(data.StateValue) > 0 {
				value = data.StateValue
			}
			return Result{
				FlagKey: flagKey,
				Value:   value,
				Reason:  ReasonRollout,
			}
		}
	}

	// No rule matched, not in rollout — return state value if set, otherwise default
	value := data.DefaultValue
	if data.StateValue != nil && len(data.StateValue) > 0 {
		value = data.StateValue
	}

	return Result{
		FlagKey: flagKey,
		Value:   value,
		Reason:  ReasonDefaultValue,
	}
}

// rolloutBucket returns a deterministic bucket in [0, 99] for a given user+flag pair.
func rolloutBucket(userID, flagKey string) int {
	h := sha256.Sum256([]byte(userID + ":" + flagKey))
	n := binary.BigEndian.Uint32(h[:4])
	return int(n % 100)
}

func matchesAllConditions(conditions []Condition, ctx EvalContext) bool {
	if len(conditions) == 0 {
		return true
	}
	for _, c := range conditions {
		attrValue, exists := lookupAttribute(c.Attribute, ctx)
		if !EvalOperator(c.Operator, attrValue, c.Value, exists) {
			return false
		}
	}
	return true
}

func lookupAttribute(attr string, ctx EvalContext) (string, bool) {
	if attr == "user_id" {
		if ctx.UserID != "" {
			return ctx.UserID, true
		}
		return "", false
	}
	v, ok := ctx.Attributes[attr]
	return v, ok
}
