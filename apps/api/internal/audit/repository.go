package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Append(ctx context.Context, input LogInput) error {
	changesJSON, err := json.Marshal(input.Changes)
	if err != nil {
		changesJSON = nil
	}

	var ipAddr *string
	if input.IPAddress != "" {
		ipAddr = &input.IPAddress
	}

	var userID *string
	if input.UserID != "" {
		userID = &input.UserID
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO audit_log (project_id, user_id, action, entity_type, entity_id, changes, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, input.ProjectID, userID, input.Action, input.EntityType, input.EntityID, changesJSON, ipAddr)
	if err != nil {
		return fmt.Errorf("appending audit log: %w", err)
	}
	return nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]Entry, int, error) {
	where := []string{"project_id = $1"}
	args := []any{params.ProjectID}
	idx := 2

	if params.UserID != "" {
		where = append(where, fmt.Sprintf("user_id = $%d", idx))
		args = append(args, params.UserID)
		idx++
	}
	if params.Action != "" {
		where = append(where, fmt.Sprintf("action = $%d", idx))
		args = append(args, params.Action)
		idx++
	}
	if params.EntityType != "" {
		where = append(where, fmt.Sprintf("entity_type = $%d", idx))
		args = append(args, params.EntityType)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM audit_log WHERE "+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting audit log: %w", err)
	}

	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	query := fmt.Sprintf(`
		SELECT id, project_id, user_id, action, entity_type, entity_id, changes, ip_address, created_at
		FROM audit_log
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, idx, idx+1)
	args = append(args, limit, params.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing audit log: %w", err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.UserID, &e.Action, &e.EntityType, &e.EntityID, &e.Changes, &e.IPAddress, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, nil
}
