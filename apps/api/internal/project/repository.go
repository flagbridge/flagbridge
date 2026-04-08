package project

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

func (r *Repository) Create(ctx context.Context, p *Project) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO projects (name, slug, description, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, p.Name, p.Slug, p.Description, p.CreatedBy).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*Project, error) {
	var p Project
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, description, created_by, created_at, updated_at
		FROM projects WHERE slug = $1
	`, slug).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) Update(ctx context.Context, slug string, req UpdateRequest) (*Project, error) {
	var p Project
	err := r.db.QueryRow(ctx, `
		UPDATE projects
		SET name = COALESCE($1, name),
		    description = COALESCE($2, description),
		    updated_at = now()
		WHERE slug = $3
		RETURNING id, name, slug, description, created_by, created_at, updated_at
	`, req.Name, req.Description, slug).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating project: %w", err)
	}
	return &p, nil
}

func (r *Repository) Delete(ctx context.Context, slug string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM projects WHERE slug = $1`, slug)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

func (r *Repository) List(ctx context.Context) ([]Project, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, slug, description, created_by, created_at, updated_at
		FROM projects ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}
