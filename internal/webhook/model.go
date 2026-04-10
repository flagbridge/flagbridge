package webhook

import "time"

// Event types supported by FlagBridge webhooks.
const (
	EventFlagCreated   = "flag.created"
	EventFlagUpdated   = "flag.updated"
	EventFlagDeleted   = "flag.deleted"
	EventFlagToggled   = "flag.toggled"
	EventFlagRollout   = "flag.rollout_changed"
	EventProjectCreated = "project.created"
	EventProjectDeleted = "project.deleted"
	EventEnvCreated    = "environment.created"
	EventKeyCreated    = "api_key.created"
)

var AllEvents = []string{
	EventFlagCreated, EventFlagUpdated, EventFlagDeleted,
	EventFlagToggled, EventFlagRollout,
	EventProjectCreated, EventProjectDeleted,
	EventEnvCreated, EventKeyCreated,
}

type Webhook struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	URL       string    `json:"url"`
	Secret    string    `json:"-"`
	Events    []string  `json:"events"`
	Active    bool      `json:"enabled"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeliveryLog struct {
	ID         string    `json:"id"`
	WebhookID  string    `json:"webhook_id"`
	EventType  string    `json:"event_type"`
	Payload    string    `json:"payload"`
	StatusCode int       `json:"status_code"`
	Response   string    `json:"response,omitempty"`
	Attempts   int       `json:"attempts"`
	Success    bool      `json:"success"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Secret string   `json:"secret,omitempty"`
}

type UpdateRequest struct {
	URL    *string  `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Active *bool    `json:"enabled,omitempty"`
}
