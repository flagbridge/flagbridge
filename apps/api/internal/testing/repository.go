package testing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSession(ctx context.Context, s *Session) error {
	overridesJSON, err := json.Marshal(s.Overrides)
	if err != nil {
		return fmt.Errorf("marshaling overrides: %w", err)
	}

	return r.db.QueryRow(ctx, `
		INSERT INTO testing_sessions (project_id, label, overrides, created_by, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, s.ProjectID, s.Label, overridesJSON, s.CreatedBy, s.ExpiresAt).Scan(&s.ID, &s.CreatedAt)
}

func (r *Repository) GetSession(ctx context.Context, id string) (*Session, error) {
	var s Session
	var overridesJSON []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, project_id, label, overrides, created_by, created_at, expires_at
		FROM testing_sessions
		WHERE id = $1 AND expires_at > now()
	`, id).Scan(&s.ID, &s.ProjectID, &s.Label, &overridesJSON, &s.CreatedBy, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}
	if err := json.Unmarshal(overridesJSON, &s.Overrides); err != nil {
		return nil, fmt.Errorf("unmarshaling overrides: %w", err)
	}
	return &s, nil
}

func (r *Repository) DeleteSession(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM testing_sessions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (r *Repository) UpdateOverrides(ctx context.Context, id string, overrides map[string]Override) error {
	overridesJSON, err := json.Marshal(overrides)
	if err != nil {
		return fmt.Errorf("marshaling overrides: %w", err)
	}

	ct, err := r.db.Exec(ctx, `
		UPDATE testing_sessions SET overrides = $1
		WHERE id = $2 AND expires_at > now()
	`, overridesJSON, id)
	if err != nil {
		return fmt.Errorf("updating overrides: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("session not found or expired")
	}
	return nil
}

func (r *Repository) ListSessions(ctx context.Context, projectID string) ([]Session, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, project_id, label, overrides, created_by, created_at, expires_at
		FROM testing_sessions
		WHERE project_id = $1 AND expires_at > now()
		ORDER BY created_at DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		var overridesJSON []byte
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Label, &overridesJSON, &s.CreatedBy, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(overridesJSON, &s.Overrides); err != nil {
			s.Overrides = make(map[string]Override)
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// CleanupExpired removes expired sessions. Call periodically.
func (r *Repository) CleanupExpired(ctx context.Context) (int64, error) {
	ct, err := r.db.Exec(ctx, `DELETE FROM testing_sessions WHERE expires_at <= now()`)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}
