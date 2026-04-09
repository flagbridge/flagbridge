package audit

import "time"

type Entry struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"project_id"`
	UserID     *string   `json:"user_id,omitempty"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Changes    any       `json:"changes,omitempty"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListParams struct {
	ProjectID  string
	UserID     string
	Action     string
	EntityType string
	Limit      int
	Offset     int
}

type LogInput struct {
	ProjectID  string
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	Changes    any
	IPAddress  string
}
