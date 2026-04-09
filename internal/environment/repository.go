package environment

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

func (r *Repository) Create(ctx context.Context, e *Environment) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO environments (project_id, name, slug, color, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, e.ProjectID, e.Name, e.Slug, e.Color, e.SortOrder).Scan(&e.ID, &e.CreatedAt)
}

func (r *Repository) ListByProject(ctx context.Context, projectID string) ([]Environment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, project_id, name, slug, color, sort_order, created_at
		FROM environments WHERE project_id = $1
		ORDER BY sort_order ASC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []Environment
	for rows.Next() {
		var e Environment
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.Name, &e.Slug, &e.Color, &e.SortOrder, &e.CreatedAt); err != nil {
			return nil, err
		}
		envs = append(envs, e)
	}
	return envs, nil
}

func (r *Repository) GetBySlug(ctx context.Context, projectID, slug string) (*Environment, error) {
	var e Environment
	err := r.db.QueryRow(ctx, `
		SELECT id, project_id, name, slug, color, sort_order, created_at
		FROM environments WHERE project_id = $1 AND slug = $2
	`, projectID, slug).Scan(&e.ID, &e.ProjectID, &e.Name, &e.Slug, &e.Color, &e.SortOrder, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}
	return &e, nil
}
