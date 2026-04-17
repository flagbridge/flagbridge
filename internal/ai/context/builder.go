package context

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Builder collects project data from Postgres and applies token budgets
// per ADR-0002. One Builder per DB pool; Build is safe for concurrent use.
type Builder struct {
	db *pgxpool.Pool
}

// NewBuilder returns a Builder backed by the given pool.
func NewBuilder(db *pgxpool.Pool) *Builder {
	return &Builder{db: db}
}

// Build returns a PromptContext for the given project as seen through the given role.
// Role determines whether product_cards are surfaced (see IncludesProductCards).
//
// Budgets applied:
//   - MaxFlags flags, sorted updated_at DESC
//   - MaxRules rules across all flags, sorted priority DESC
//   - MaxAuditEntries audit entries, sorted created_at DESC
//
// When any budget is hit, Truncation is populated so the LLM and UI know the context
// is partial.
func (b *Builder) Build(ctx context.Context, projectID, role string) (PromptContext, error) {
	pc := PromptContext{Version: CurrentVersion, Role: role}

	if err := b.loadProject(ctx, projectID, &pc); err != nil {
		return pc, err
	}

	if err := b.loadFlags(ctx, projectID, &pc); err != nil {
		return pc, err
	}

	if err := b.loadRules(ctx, projectID, &pc); err != nil {
		return pc, err
	}

	if IncludesProductCards(role) {
		if err := b.loadProductCards(ctx, projectID, &pc); err != nil {
			return pc, err
		}
	}

	if err := b.loadRecentAudit(ctx, projectID, &pc); err != nil {
		return pc, err
	}

	return pc, nil
}

func (b *Builder) loadProject(ctx context.Context, projectID string, pc *PromptContext) error {
	var desc *string
	err := b.db.QueryRow(ctx, `
		SELECT id, name, slug, description
		FROM projects WHERE id = $1
	`, projectID).Scan(&pc.Project.ID, &pc.Project.Name, &pc.Project.Slug, &desc)
	if err != nil {
		return fmt.Errorf("loading project: %w", err)
	}
	if desc != nil {
		pc.Project.Description = *desc
	}
	return nil
}

func (b *Builder) loadFlags(ctx context.Context, projectID string, pc *PromptContext) error {
	var total int
	if err := b.db.QueryRow(ctx, `SELECT COUNT(*) FROM flags WHERE project_id = $1`, projectID).Scan(&total); err != nil {
		return fmt.Errorf("counting flags: %w", err)
	}

	rows, err := b.db.Query(ctx, `
		SELECT key, name, COALESCE(description, ''), type, default_value, tags, updated_at
		FROM flags
		WHERE project_id = $1
		ORDER BY updated_at DESC
		LIMIT $2
	`, projectID, MaxFlags)
	if err != nil {
		return fmt.Errorf("loading flags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f FlagContext
		var defaultValue json.RawMessage
		if err := rows.Scan(&f.Key, &f.Name, &f.Description, &f.Type, &defaultValue, &f.Tags, &f.UpdatedAt); err != nil {
			return fmt.Errorf("scanning flag: %w", err)
		}
		f.DefaultValue = string(defaultValue)
		pc.Flags = append(pc.Flags, f)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating flags: %w", err)
	}

	if total > MaxFlags {
		ensureTruncation(pc)
		pc.Truncation.FlagsTruncated = true
		pc.Truncation.TotalFlags = total
	}
	return nil
}

func (b *Builder) loadRules(ctx context.Context, projectID string, pc *PromptContext) error {
	var total int
	if err := b.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM targeting_rules r
		JOIN flags f ON f.id = r.flag_id
		WHERE f.project_id = $1
	`, projectID).Scan(&total); err != nil {
		return fmt.Errorf("counting rules: %w", err)
	}

	rows, err := b.db.Query(ctx, `
		SELECT f.key, e.slug, r.name, r.priority, r.enabled, r.conditions, r.value
		FROM targeting_rules r
		JOIN flags f ON f.id = r.flag_id
		JOIN environments e ON e.id = r.environment_id
		WHERE f.project_id = $1
		ORDER BY r.priority DESC, r.updated_at DESC
		LIMIT $2
	`, projectID, MaxRules)
	if err != nil {
		return fmt.Errorf("loading rules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r RuleContext
		var conditions, value json.RawMessage
		if err := rows.Scan(&r.FlagKey, &r.Environment, &r.Name, &r.Priority, &r.Enabled, &conditions, &value); err != nil {
			return fmt.Errorf("scanning rule: %w", err)
		}
		r.Conditions = string(conditions)
		r.Value = string(value)
		pc.Rules = append(pc.Rules, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating rules: %w", err)
	}

	if total > MaxRules {
		ensureTruncation(pc)
		pc.Truncation.RulesTruncated = true
		pc.Truncation.TotalRules = total
	}
	return nil
}

func (b *Builder) loadProductCards(ctx context.Context, projectID string, pc *PromptContext) error {
	rows, err := b.db.Query(ctx, `
		SELECT f.key, pc.hypothesis, pc.success_metrics, pc.go_no_go, pc.status
		FROM product_cards pc
		JOIN flags f ON f.id = pc.flag_id
		WHERE pc.project_id = $1
		ORDER BY pc.updated_at DESC
	`, projectID)
	if err != nil {
		return fmt.Errorf("loading product cards: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c ProductCardContext
		if err := rows.Scan(&c.FlagKey, &c.Hypothesis, &c.SuccessMetrics, &c.GoNoGo, &c.Status); err != nil {
			return fmt.Errorf("scanning product card: %w", err)
		}
		pc.ProductCards = append(pc.ProductCards, c)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating product cards: %w", err)
	}
	return nil
}

func (b *Builder) loadRecentAudit(ctx context.Context, projectID string, pc *PromptContext) error {
	var total int
	if err := b.db.QueryRow(ctx, `SELECT COUNT(*) FROM audit_log WHERE project_id = $1`, projectID).Scan(&total); err != nil {
		return fmt.Errorf("counting audit entries: %w", err)
	}

	rows, err := b.db.Query(ctx, `
		SELECT a.action, a.entity_type, a.entity_id, COALESCE(u.name, ''), a.created_at
		FROM audit_log a
		LEFT JOIN users u ON u.id = a.user_id
		WHERE a.project_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2
	`, projectID, MaxAuditEntries)
	if err != nil {
		return fmt.Errorf("loading audit log: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a AuditEntry
		if err := rows.Scan(&a.Action, &a.EntityType, &a.EntityID, &a.ActorName, &a.CreatedAt); err != nil {
			return fmt.Errorf("scanning audit entry: %w", err)
		}
		pc.RecentAudit = append(pc.RecentAudit, a)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating audit: %w", err)
	}

	if total > MaxAuditEntries {
		ensureTruncation(pc)
		pc.Truncation.AuditEntriesTruncated = true
		pc.Truncation.TotalAuditEntries = total
	}
	return nil
}

func ensureTruncation(pc *PromptContext) {
	if pc.Truncation == nil {
		pc.Truncation = &TruncationInfo{}
	}
}
