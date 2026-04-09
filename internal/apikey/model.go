package apikey

import "time"

type APIKey struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	KeyPrefix     string     `json:"key_prefix"`
	Scope         string     `json:"scope"`
	ProjectID     string     `json:"project_id"`
	EnvironmentID string     `json:"environment_id,omitempty"`
	CreatedBy     string     `json:"created_by"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateRequest struct {
	Name          string `json:"name"`
	Scope         string `json:"scope"`
	ProjectID     string `json:"project_id"`
	EnvironmentID string `json:"environment_id,omitempty"`
}

type CreateResponse struct {
	APIKey
	FullKey string `json:"key"`
}
