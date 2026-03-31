package flag

import (
	"encoding/json"
	"time"
)

type Flag struct {
	ID           string          `json:"id"`
	ProjectID    string          `json:"project_id"`
	Key          string          `json:"key"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Type         string          `json:"type"`
	DefaultValue json.RawMessage `json:"default_value"`
	Tags         []string        `json:"tags"`
	CreatedBy    string          `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type FlagState struct {
	ID            string          `json:"id"`
	FlagID        string          `json:"flag_id"`
	EnvironmentID string          `json:"environment_id"`
	Enabled       bool            `json:"enabled"`
	Value         json.RawMessage `json:"value,omitempty"`
	UpdatedBy     string          `json:"updated_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type CreateRequest struct {
	Key          string          `json:"key"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Type         string          `json:"type"`
	DefaultValue json.RawMessage `json:"default_value,omitempty"`
	Tags         []string        `json:"tags,omitempty"`
}

type UpdateRequest struct {
	Name         *string          `json:"name,omitempty"`
	Description  *string          `json:"description,omitempty"`
	DefaultValue json.RawMessage  `json:"default_value,omitempty"`
	Tags         []string         `json:"tags,omitempty"`
}

type SetStateRequest struct {
	Enabled bool            `json:"enabled"`
	Value   json.RawMessage `json:"value,omitempty"`
}
