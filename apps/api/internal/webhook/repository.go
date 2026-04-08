package webhook

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

func (r *Repository) Create(ctx context.Context, w *Webhook) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO webhooks (project_id, url, secret, events, active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, w.ProjectID, w.URL, w.Secret, w.Events, w.Active, w.CreatedBy).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Webhook, error) {
	var w Webhook
	err := r.db.QueryRow(ctx, `
		SELECT id, project_id, url, secret, events, active, created_by, created_at, updated_at
		FROM webhooks WHERE id = $1
	`, id).Scan(&w.ID, &w.ProjectID, &w.URL, &w.Secret, &w.Events, &w.Active, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}
	return &w, nil
}

func (r *Repository) ListByProject(ctx context.Context, projectID string) ([]Webhook, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, project_id, url, events, active, created_by, created_at, updated_at
		FROM webhooks WHERE project_id = $1 ORDER BY created_at DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []Webhook
	for rows.Next() {
		var w Webhook
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &w.Events, &w.Active, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateRequest) (*Webhook, error) {
	var w Webhook
	err := r.db.QueryRow(ctx, `
		UPDATE webhooks
		SET url = COALESCE($1, url),
		    events = COALESCE($2, events),
		    active = COALESCE($3, active),
		    updated_at = now()
		WHERE id = $4
		RETURNING id, project_id, url, events, active, created_by, created_at, updated_at
	`, req.URL, req.Events, req.Active, id).Scan(&w.ID, &w.ProjectID, &w.URL, &w.Events, &w.Active, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating webhook: %w", err)
	}
	return &w, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM webhooks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting webhook: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("webhook not found")
	}
	return nil
}

func (r *Repository) LogDelivery(ctx context.Context, log *DeliveryLog) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO webhook_delivery_logs (webhook_id, event_type, payload, status_code, response, attempts, success)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`, log.WebhookID, log.EventType, log.Payload, log.StatusCode, log.Response, log.Attempts, log.Success).Scan(&log.ID, &log.CreatedAt)
}

func (r *Repository) ListDeliveryLogs(ctx context.Context, webhookID string, limit int) ([]DeliveryLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, webhook_id, event_type, payload, status_code, response, attempts, success, created_at
		FROM webhook_delivery_logs
		WHERE webhook_id = $1
		ORDER BY created_at DESC LIMIT $2
	`, webhookID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []DeliveryLog
	for rows.Next() {
		var l DeliveryLog
		if err := rows.Scan(&l.ID, &l.WebhookID, &l.EventType, &l.Payload, &l.StatusCode, &l.Response, &l.Attempts, &l.Success, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// FindByEvent returns all active webhooks for a project subscribed to a given event type.
func (r *Repository) FindByEvent(ctx context.Context, projectID, eventType string) ([]Webhook, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, project_id, url, secret, events, active, created_by, created_at, updated_at
		FROM webhooks
		WHERE project_id = $1 AND active = true AND $2 = ANY(events)
	`, projectID, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []Webhook
	for rows.Next() {
		var w Webhook
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &w.Secret, &w.Events, &w.Active, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, nil
}
