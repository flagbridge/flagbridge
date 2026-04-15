package member

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByProject(ctx context.Context, projectID string) ([]Member, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role, u.email, u.name,
		       pm.created_at, pm.updated_at
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at ASC
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing members: %w", err)
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role,
			&m.Email, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning member: %w", err)
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *Repository) Add(ctx context.Context, projectID, userID, role string) (*Member, error) {
	var m Member
	err := r.db.QueryRow(ctx, `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, user_id) DO UPDATE SET role = $3, updated_at = now()
		RETURNING id, project_id, user_id, role, created_at, updated_at
	`, projectID, userID, role).Scan(
		&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("adding member: %w", err)
	}

	// Fetch user info
	_ = r.db.QueryRow(ctx, `SELECT email, name FROM users WHERE id = $1`, userID).Scan(&m.Email, &m.Name)
	return &m, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (userID string, err error) {
	err = r.db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	return userID, nil
}

func (r *Repository) UpdateRole(ctx context.Context, memberID, role string) (*Member, error) {
	var m Member
	err := r.db.QueryRow(ctx, `
		UPDATE project_members SET role = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, project_id, user_id, role, created_at, updated_at
	`, memberID, role).Scan(
		&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("updating member role: %w", err)
	}

	_ = r.db.QueryRow(context.Background(), `SELECT email, name FROM users WHERE id = $1`, m.UserID).Scan(&m.Email, &m.Name)
	return &m, nil
}

func (r *Repository) Remove(ctx context.Context, memberID string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM project_members WHERE id = $1`, memberID)
	if err != nil {
		return fmt.Errorf("removing member: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}
