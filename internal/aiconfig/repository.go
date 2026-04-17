package aiconfig

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

func (r *Repository) Upsert(ctx context.Context, projectID, provider, modelID, encryptedKey string, baseURL *string, maxTokens int, temperature float64) (*Provider, error) {
	var p Provider
	var rawEncKey *string

	err := r.db.QueryRow(ctx, `
		INSERT INTO ai_providers (project_id, provider, model_id, encrypted_api_key, base_url, max_tokens, temperature)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (project_id) DO UPDATE SET
			provider = EXCLUDED.provider,
			model_id = EXCLUDED.model_id,
			encrypted_api_key = COALESCE(EXCLUDED.encrypted_api_key, ai_providers.encrypted_api_key),
			base_url = EXCLUDED.base_url,
			max_tokens = EXCLUDED.max_tokens,
			temperature = EXCLUDED.temperature,
			updated_at = now()
		RETURNING id, project_id, provider, model_id, encrypted_api_key, base_url,
		          max_tokens, temperature, monthly_usage, monthly_limit, last_reset_at,
		          created_at, updated_at
	`, projectID, provider, modelID, nilIfEmpty(encryptedKey), baseURL, maxTokens, temperature).Scan(
		&p.ID, &p.ProjectID, &p.ProviderName, &p.ModelID, &rawEncKey, &p.BaseURL,
		&p.MaxTokens, &p.Temperature, &p.MonthlyUsage, &p.MonthlyLimit, &p.LastResetAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upserting AI provider: %w", err)
	}

	p.HasAPIKey = rawEncKey != nil && *rawEncKey != ""
	return &p, nil
}

func (r *Repository) GetByProjectID(ctx context.Context, projectID string) (*Provider, error) {
	var p Provider
	var rawEncKey *string

	err := r.db.QueryRow(ctx, `
		SELECT id, project_id, provider, model_id, encrypted_api_key, base_url,
		       max_tokens, temperature, monthly_usage, monthly_limit, last_reset_at,
		       created_at, updated_at
		FROM ai_providers WHERE project_id = $1
	`, projectID).Scan(
		&p.ID, &p.ProjectID, &p.ProviderName, &p.ModelID, &rawEncKey, &p.BaseURL,
		&p.MaxTokens, &p.Temperature, &p.MonthlyUsage, &p.MonthlyLimit, &p.LastResetAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("fetching AI provider: %w", err)
	}

	p.HasAPIKey = rawEncKey != nil && *rawEncKey != ""
	return &p, nil
}

func (r *Repository) GetEncryptedKey(ctx context.Context, projectID string) (string, error) {
	var encKey *string
	err := r.db.QueryRow(ctx, `
		SELECT encrypted_api_key FROM ai_providers WHERE project_id = $1
	`, projectID).Scan(&encKey)
	if err != nil {
		return "", fmt.Errorf("fetching encrypted key: %w", err)
	}
	if encKey == nil {
		return "", fmt.Errorf("no API key configured")
	}
	return *encKey, nil
}

func (r *Repository) Delete(ctx context.Context, projectID string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM ai_providers WHERE project_id = $1`, projectID)
	if err != nil {
		return fmt.Errorf("deleting AI provider: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("AI provider config not found")
	}
	return nil
}

func (r *Repository) IncrementUsage(ctx context.Context, projectID string) (currentUsage int, limit int, err error) {
	err = r.db.QueryRow(ctx, `
		UPDATE ai_providers
		SET monthly_usage = CASE
			WHEN last_reset_at < date_trunc('month', now())
			THEN 1
			ELSE monthly_usage + 1
		END,
		last_reset_at = CASE
			WHEN last_reset_at < date_trunc('month', now())
			THEN date_trunc('month', now())
			ELSE last_reset_at
		END,
		updated_at = now()
		WHERE project_id = $1
		RETURNING monthly_usage, monthly_limit
	`, projectID).Scan(&currentUsage, &limit)
	if err != nil {
		return 0, 0, fmt.Errorf("incrementing usage: %w", err)
	}
	return currentUsage, limit, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
