package productcard

import "time"

// ProductCard represents the product context attached to a feature flag.
// Every flag can have one product card with hypothesis, metrics, and lifecycle status.
type ProductCard struct {
	ID             string    `json:"id"`
	FlagID         string    `json:"flag_id"`
	ProjectID      string    `json:"project_id"`
	Hypothesis     string    `json:"hypothesis"`
	SuccessMetrics string    `json:"success_metrics"`
	GoNoGo         string    `json:"go_no_go"`
	OwnerID        *string   `json:"owner_id,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UpsertRequest is used for PUT /projects/:slug/flags/:key/product-card.
// All fields are optional — missing fields retain their current value on update.
type UpsertRequest struct {
	Hypothesis     *string `json:"hypothesis,omitempty"`
	SuccessMetrics *string `json:"success_metrics,omitempty"`
	GoNoGo         *string `json:"go_no_go,omitempty"`
	OwnerID        *string `json:"owner_id,omitempty"`
	Status         *string `json:"status,omitempty"`
}

// ValidStatuses lists the allowed lifecycle statuses for a product card.
var ValidStatuses = map[string]bool{
	"planning":   true,
	"active":     true,
	"rolled_out": true,
	"archived":   true,
}
