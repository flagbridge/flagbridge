package flag

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

func (r *Repository) Create(ctx context.Context, f *Flag) error {
	tags := f.Tags
	if tags == nil {
		tags = []string{}
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO flags (project_id, key, name, description, type, default_value, tags, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, f.ProjectID, f.Key, f.Name, f.Description, f.Type, f.DefaultValue, tags, f.CreatedBy).Scan(
		&f.ID, &f.CreatedAt, &f.UpdatedAt,
	)
}

func (r *Repository) GetByKey(ctx context.Context, projectID, key string) (*Flag, error) {
	var f Flag
	var tags []string
	err := r.db.QueryRow(ctx, `
		SELECT id, project_id, key, name, description, type, default_value, tags, created_by, created_at, updated_at
		FROM flags WHERE project_id = $1 AND key = $2
	`, projectID, key).Scan(
		&f.ID, &f.ProjectID, &f.Key, &f.Name, &f.Description, &f.Type, &f.DefaultValue, &tags, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("flag not found: %w", err)
	}
	f.Tags = tags
	return &f, nil
}

func (r *Repository) ListByProject(ctx context.Context, projectID string) ([]Flag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, project_id, key, name, description, type, default_value, tags, created_by, created_at, updated_at
		FROM flags WHERE project_id = $1 ORDER BY created_at DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []Flag
	for rows.Next() {
		var f Flag
		var tags []string
		if err := rows.Scan(
			&f.ID, &f.ProjectID, &f.Key, &f.Name, &f.Description, &f.Type, &f.DefaultValue, &tags, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		f.Tags = tags
		flags = append(flags, f)
	}
	return flags, nil
}

func (r *Repository) Update(ctx context.Context, projectID, key string, req UpdateRequest) (*Flag, error) {
	var f Flag
	var tags []string

	defaultValue := json.RawMessage("null")
	if req.DefaultValue != nil {
		defaultValue = req.DefaultValue
	}

	err := r.db.QueryRow(ctx, `
		UPDATE flags SET
			name = COALESCE($3, name),
			description = COALESCE($4, description),
			default_value = CASE WHEN $5::jsonb IS NOT NULL THEN $5::jsonb ELSE default_value END,
			tags = COALESCE($6, tags),
			updated_at = now()
		WHERE project_id = $1 AND key = $2
		RETURNING id, project_id, key, name, description, type, default_value, tags, created_by, created_at, updated_at
	`, projectID, key, req.Name, req.Description, defaultValue, req.Tags).Scan(
		&f.ID, &f.ProjectID, &f.Key, &f.Name, &f.Description, &f.Type, &f.DefaultValue, &tags, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("updating flag: %w", err)
	}
	f.Tags = tags
	return &f, nil
}

func (r *Repository) Delete(ctx context.Context, projectID, key string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM flags WHERE project_id = $1 AND key = $2`, projectID, key)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("flag not found")
	}
	return nil
}

func (r *Repository) GetState(ctx context.Context, flagID, envID string) (*FlagState, error) {
	var fs FlagState
	err := r.db.QueryRow(ctx, `
		SELECT id, flag_id, environment_id, enabled, value, updated_by, updated_at
		FROM flag_states WHERE flag_id = $1 AND environment_id = $2
	`, flagID, envID).Scan(
		&fs.ID, &fs.FlagID, &fs.EnvironmentID, &fs.Enabled, &fs.Value, &fs.UpdatedBy, &fs.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("flag state not found: %w", err)
	}
	return &fs, nil
}

func (r *Repository) SetState(ctx context.Context, flagID, envID, userID string, req SetStateRequest) (*FlagState, error) {
	var fs FlagState
	err := r.db.QueryRow(ctx, `
		INSERT INTO flag_states (flag_id, environment_id, enabled, value, updated_by)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (flag_id, environment_id)
		DO UPDATE SET enabled = $3, value = $4, updated_by = $5, updated_at = now()
		RETURNING id, flag_id, environment_id, enabled, value, updated_by, updated_at
	`, flagID, envID, req.Enabled, req.Value, userID).Scan(
		&fs.ID, &fs.FlagID, &fs.EnvironmentID, &fs.Enabled, &fs.Value, &fs.UpdatedBy, &fs.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("setting flag state: %w", err)
	}
	return &fs, nil
}
