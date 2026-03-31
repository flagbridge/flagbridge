package targeting

import (
	"encoding/json"
	"time"

	"github.com/flagbridge/flagbridge/apps/api/internal/evaluation"
)

type Rule struct {
	ID            string                 `json:"id"`
	FlagID        string                 `json:"flag_id"`
	EnvironmentID string                 `json:"environment_id"`
	Name          string                 `json:"name"`
	Priority      int                    `json:"priority"`
	Conditions    []evaluation.Condition `json:"conditions"`
	Value         json.RawMessage        `json:"value"`
	Enabled       bool                   `json:"enabled"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type SetRulesRequest struct {
	Rules []RuleInput `json:"rules"`
}

type RuleInput struct {
	Name       string                 `json:"name"`
	Priority   int                    `json:"priority"`
	Conditions []evaluation.Condition `json:"conditions"`
	Value      json.RawMessage        `json:"value"`
	Enabled    bool                   `json:"enabled"`
}
