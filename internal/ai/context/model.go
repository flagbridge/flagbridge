// Package context builds versioned, token-budgeted prompt contexts
// for LLM interactions per ADR-0002 (docs/adr/0002-ai-prompt-context-struct.md).
//
// PromptContext is serialized to XML (Anthropic) or JSON (OpenAI) via Render*.
// A Builder collects data from Postgres and applies truncation budgets.
package context

import "time"

// CurrentVersion of the PromptContext schema. Bump on breaking changes —
// clients pin on major to keep prompt behaviour stable across server upgrades.
const CurrentVersion = "1.0"

// Budget caps. Static in V1 — see ADR-0002 "Future" section for adaptive heuristics.
const (
	// MaxFlags is the most flags included per context, sorted by updated_at DESC.
	MaxFlags = 50
	// MaxAuditEntries is the most audit events included, sorted by created_at DESC.
	MaxAuditEntries = 20
	// MaxRules is the most targeting rules included across all flags in the project,
	// sorted by priority DESC.
	MaxRules = 100
)

// Role values as understood by the context builder. The actual role string comes
// from the authenticated session; these constants are used for renderer branching.
const (
	RoleEngineer = "engineer"
	RoleProduct  = "product"
	RoleAdmin    = "admin"
	RoleViewer   = "viewer"
)

// PromptContext is the canonical shape handed to any LLM provider via the renderer.
type PromptContext struct {
	Version      string               `json:"version"`
	Project      ProjectContext       `json:"project"`
	Flags        []FlagContext        `json:"flags"`
	Rules        []RuleContext        `json:"rules,omitempty"`
	ProductCards []ProductCardContext `json:"product_cards,omitempty"`
	RecentAudit  []AuditEntry         `json:"recent_audit,omitempty"`
	Role         string               `json:"role"`
	Truncation   *TruncationInfo      `json:"truncation,omitempty"`
}

// ProjectContext carries the minimum project identity. Internal IDs are included so
// the LLM can disambiguate when a user references the project by name in prompts.
type ProjectContext struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
}

// FlagContext is a flattened flag record. DefaultValue is kept as raw JSON string so
// the LLM sees booleans, numbers and objects in their natural form.
type FlagContext struct {
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Type         string    `json:"type"`
	DefaultValue string    `json:"default_value"`
	Tags         []string  `json:"tags,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RuleContext groups a targeting rule by flag key + environment slug to avoid forcing
// the LLM to resolve UUID references.
type RuleContext struct {
	FlagKey     string `json:"flag_key"`
	Environment string `json:"environment"`
	Name        string `json:"name,omitempty"`
	Priority    int    `json:"priority"`
	Enabled     bool   `json:"enabled"`
	Conditions  string `json:"conditions"`
	Value       string `json:"value"`
}

// ProductCardContext is included only when Role == RoleProduct or RoleAdmin, per ADR-0002.
// Engineers don't need product intelligence data surfaced in their AI prompts.
type ProductCardContext struct {
	FlagKey        string `json:"flag_key"`
	Hypothesis     string `json:"hypothesis,omitempty"`
	SuccessMetrics string `json:"success_metrics,omitempty"`
	GoNoGo         string `json:"go_no_go,omitempty"`
	Status         string `json:"status"`
}

// AuditEntry is a sanitised audit record — no IP, no changes blob. LLM doesn't need PII.
type AuditEntry struct {
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	ActorName  string    `json:"actor_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// TruncationInfo is non-nil when any budget was applied. Clients (UI, tests) can show
// "showing 50 of N flags" badges; LLM sees it and knows context is partial.
type TruncationInfo struct {
	FlagsTruncated        bool `json:"flags_truncated,omitempty"`
	TotalFlags            int  `json:"total_flags,omitempty"`
	RulesTruncated        bool `json:"rules_truncated,omitempty"`
	TotalRules            int  `json:"total_rules,omitempty"`
	AuditEntriesTruncated bool `json:"audit_entries_truncated,omitempty"`
	TotalAuditEntries     int  `json:"total_audit_entries,omitempty"`
}

// IncludesProductCards returns true when the role warrants surfacing product intelligence.
func IncludesProductCards(role string) bool {
	return role == RoleProduct || role == RoleAdmin
}
