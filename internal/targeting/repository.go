package targeting

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/flagbridge/flagbridge/internal/evaluation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByFlagAndEnv(ctx context.Context, flagID, envID string) ([]Rule, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, flag_id, environment_id, name, priority, conditions, value, enabled, created_at, updated_at
		FROM targeting_rules
		WHERE flag_id = $1 AND environment_id = $2
		ORDER BY priority ASC
	`, flagID, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		var condJSON []byte
		if err := rows.Scan(
			&rule.ID, &rule.FlagID, &rule.EnvironmentID, &rule.Name,
			&rule.Priority, &condJSON, &rule.Value, &rule.Enabled,
			&rule.CreatedAt, &rule.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(condJSON, &rule.Conditions); err != nil {
			rule.Conditions = []evaluation.Condition{}
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *Repository) SetRules(ctx context.Context, flagID, envID string, inputs []RuleInput) ([]Rule, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete existing rules
	if _, err := tx.Exec(ctx, `DELETE FROM targeting_rules WHERE flag_id = $1 AND environment_id = $2`, flagID, envID); err != nil {
		return nil, fmt.Errorf("deleting old rules: %w", err)
	}

	var rules []Rule
	for _, input := range inputs {
		condJSON, err := json.Marshal(input.Conditions)
		if err != nil {
			return nil, fmt.Errorf("marshaling conditions: %w", err)
		}

		var rule Rule
		var condBytes []byte
		err = tx.QueryRow(ctx, `
			INSERT INTO targeting_rules (flag_id, environment_id, name, priority, conditions, value, enabled)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, flag_id, environment_id, name, priority, conditions, value, enabled, created_at, updated_at
		`, flagID, envID, input.Name, input.Priority, condJSON, input.Value, input.Enabled).Scan(
			&rule.ID, &rule.FlagID, &rule.EnvironmentID, &rule.Name,
			&rule.Priority, &condBytes, &rule.Value, &rule.Enabled,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("inserting rule: %w", err)
		}
		json.Unmarshal(condBytes, &rule.Conditions)
		rules = append(rules, rule)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return rules, nil
}
