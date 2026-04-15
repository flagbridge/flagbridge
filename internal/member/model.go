package member

import "time"

// Member represents a user's membership in a project with a specific role.
type Member struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

// ValidRoles lists the 4 project-level roles.
var ValidRoles = map[string]bool{
	"admin":    true,
	"engineer": true,
	"product":  true,
	"viewer":   true,
}
