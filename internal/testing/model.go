package testing

import (
	"encoding/json"
	"time"
)

type Session struct {
	ID        string            `json:"id"`
	ProjectID string            `json:"project_id"`
	Label     string            `json:"label,omitempty"`
	Overrides map[string]Override `json:"overrides"`
	CreatedBy string            `json:"created_by"`
	CreatedAt time.Time         `json:"created_at"`
	ExpiresAt time.Time         `json:"expires_at"`
}

type Override struct {
	FlagKey     string          `json:"flag_key"`
	Value       json.RawMessage `json:"value"`
	Environment string          `json:"environment,omitempty"`
}

type CreateSessionRequest struct {
	ProjectID string `json:"project_id"`
	Label     string `json:"label,omitempty"`
	TTLSecs   int    `json:"ttl_secs,omitempty"` // default 3600 (1 hour)
}

type SetOverrideRequest struct {
	FlagKey     string          `json:"flag_key"`
	Value       json.RawMessage `json:"value"`
	Environment string          `json:"environment,omitempty"`
}

type SetOverridesBatchRequest struct {
	Overrides []SetOverrideRequest `json:"overrides"`
}
