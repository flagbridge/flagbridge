package apikey

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

func (r *Repository) Create(ctx context.Context, name, keyHash, keyPrefix, scope, projectID string, environmentID *string, createdBy string) (*APIKey, error) {
	var ak APIKey

	err := r.db.QueryRow(ctx, `
		INSERT INTO api_keys (name, key_hash, key_prefix, scope, project_id, environment_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, key_prefix, scope, project_id, environment_id, created_by, created_at
	`, name, keyHash, keyPrefix, scope, projectID, environmentID, createdBy).Scan(
		&ak.ID, &ak.Name, &ak.KeyPrefix, &ak.Scope, &ak.ProjectID, &ak.EnvironmentID, &ak.CreatedBy, &ak.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating API key: %w", err)
	}
	return &ak, nil
}

func (r *Repository) List(ctx context.Context) ([]APIKey, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, key_prefix, scope, project_id, environment_id, created_by, last_used_at, expires_at, created_at
		FROM api_keys ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var ak APIKey
		if err := rows.Scan(
			&ak.ID, &ak.Name, &ak.KeyPrefix, &ak.Scope, &ak.ProjectID,
			&ak.EnvironmentID, &ak.CreatedBy, &ak.LastUsedAt, &ak.ExpiresAt, &ak.CreatedAt,
		); err != nil {
			return nil, err
		}
		keys = append(keys, ak)
	}
	return keys, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM api_keys WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}
	return nil
}
