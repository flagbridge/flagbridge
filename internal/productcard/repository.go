package productcard

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByFlagID(ctx context.Context, flagID string) (*ProductCard, error) {
	var c ProductCard
	err := r.db.QueryRow(ctx, `
		SELECT id, flag_id, project_id, hypothesis, success_metrics, go_no_go,
		       owner_id, status, created_at, updated_at
		FROM product_cards WHERE flag_id = $1
	`, flagID).Scan(
		&c.ID, &c.FlagID, &c.ProjectID, &c.Hypothesis, &c.SuccessMetrics,
		&c.GoNoGo, &c.OwnerID, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting product card: %w", err)
	}
	return &c, nil
}

func (r *Repository) Upsert(ctx context.Context, flagID, projectID string, req UpsertRequest) (*ProductCard, error) {
	var c ProductCard
	err := r.db.QueryRow(ctx, `
		INSERT INTO product_cards (flag_id, project_id, hypothesis, success_metrics, go_no_go, owner_id, status)
		VALUES ($1, $2, COALESCE($3, ''), COALESCE($4, ''), COALESCE($5, ''), $6, COALESCE($7, 'planning'))
		ON CONFLICT (flag_id)
		DO UPDATE SET
			hypothesis     = COALESCE($3, product_cards.hypothesis),
			success_metrics = COALESCE($4, product_cards.success_metrics),
			go_no_go       = COALESCE($5, product_cards.go_no_go),
			owner_id       = CASE WHEN $6 IS NOT NULL THEN $6 ELSE product_cards.owner_id END,
			status         = COALESCE($7, product_cards.status),
			updated_at     = now()
		RETURNING id, flag_id, project_id, hypothesis, success_metrics, go_no_go,
		          owner_id, status, created_at, updated_at
	`, flagID, projectID,
		req.Hypothesis, req.SuccessMetrics, req.GoNoGo, req.OwnerID, req.Status,
	).Scan(
		&c.ID, &c.FlagID, &c.ProjectID, &c.Hypothesis, &c.SuccessMetrics,
		&c.GoNoGo, &c.OwnerID, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upserting product card: %w", err)
	}
	return &c, nil
}

func (r *Repository) Delete(ctx context.Context, flagID string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM product_cards WHERE flag_id = $1`, flagID)
	if err != nil {
		return fmt.Errorf("deleting product card: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("product card not found")
	}
	return nil
}
