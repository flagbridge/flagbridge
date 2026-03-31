package environment

import "time"

type Environment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Color     string    `json:"color,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color,omitempty"`
	SortOrder int    `json:"sort_order"`
}
