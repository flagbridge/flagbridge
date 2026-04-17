package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ContextBuilder collects project data to build a rich system prompt for LLM interactions.
type ContextBuilder struct {
	db *pgxpool.Pool
}

func NewContextBuilder(db *pgxpool.Pool) *ContextBuilder {
	return &ContextBuilder{db: db}
}

type projectContext struct {
	ProjectName string         `json:"project_name"`
	Flags       []flagContext  `json:"flags"`
	Members     []memberBrief  `json:"members"`
	RecentAudit []auditBrief   `json:"recent_audit"`
}

type flagContext struct {
	Key          string          `json:"key"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Type         string          `json:"type"`
	DefaultValue json.RawMessage `json:"default_value"`
	Tags         []string        `json:"tags,omitempty"`
	Hypothesis   string          `json:"hypothesis,omitempty"`
	Metrics      string          `json:"success_metrics,omitempty"`
	Status       string          `json:"card_status,omitempty"`
	RuleCount    int             `json:"rule_count"`
}

type memberBrief struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type auditBrief struct {
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// Build collects all relevant project data and returns a structured context string.
func (cb *ContextBuilder) Build(ctx context.Context, projectID string) (string, error) {
	pc := projectContext{}

	// Project name
	err := cb.db.QueryRow(ctx, `SELECT name FROM projects WHERE id = $1`, projectID).Scan(&pc.ProjectName)
	if err != nil {
		return "", fmt.Errorf("loading project: %w", err)
	}

	// Flags + product cards + rule counts (single query)
	rows, err := cb.db.Query(ctx, `
		SELECT f.key, f.name, f.description, f.type, f.default_value, f.tags,
		       COALESCE(pc.hypothesis, ''), COALESCE(pc.success_metrics, ''), COALESCE(pc.status, ''),
		       COALESCE(rc.rule_count, 0)
		FROM flags f
		LEFT JOIN product_cards pc ON pc.flag_id = f.id
		LEFT JOIN (
			SELECT flag_id, COUNT(*) AS rule_count
			FROM targeting_rules
			GROUP BY flag_id
		) rc ON rc.flag_id = f.id
		WHERE f.project_id = $1
		ORDER BY f.created_at DESC
		LIMIT 50
	`, projectID)
	if err != nil {
		return "", fmt.Errorf("loading flags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fc flagContext
		if err := rows.Scan(
			&fc.Key, &fc.Name, &fc.Description, &fc.Type, &fc.DefaultValue, &fc.Tags,
			&fc.Hypothesis, &fc.Metrics, &fc.Status, &fc.RuleCount,
		); err != nil {
			return "", fmt.Errorf("scanning flag: %w", err)
		}
		pc.Flags = append(pc.Flags, fc)
	}

	// Members
	memberRows, err := cb.db.Query(ctx, `
		SELECT u.name, pm.role
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.role, u.name
	`, projectID)
	if err != nil {
		return "", fmt.Errorf("loading members: %w", err)
	}
	defer memberRows.Close()

	for memberRows.Next() {
		var m memberBrief
		if err := memberRows.Scan(&m.Name, &m.Role); err != nil {
			return "", fmt.Errorf("scanning member: %w", err)
		}
		pc.Members = append(pc.Members, m)
	}

	// Recent audit (last 20 entries)
	auditRows, err := cb.db.Query(ctx, `
		SELECT action, entity_type, entity_id
		FROM audit_log
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT 20
	`, projectID)
	if err != nil {
		return "", fmt.Errorf("loading audit: %w", err)
	}
	defer auditRows.Close()

	for auditRows.Next() {
		var a auditBrief
		if err := auditRows.Scan(&a.Action, &a.EntityType, &a.EntityID); err != nil {
			return "", fmt.Errorf("scanning audit: %w", err)
		}
		pc.RecentAudit = append(pc.RecentAudit, a)
	}

	return formatContext(pc), nil
}

func formatContext(pc projectContext) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Project: %s\n\n", pc.ProjectName))

	sb.WriteString(fmt.Sprintf("## Flags (%d)\n", len(pc.Flags)))
	for _, f := range pc.Flags {
		sb.WriteString(fmt.Sprintf("- **%s** (%s, %s)", f.Key, f.Name, f.Type))
		if f.RuleCount > 0 {
			sb.WriteString(fmt.Sprintf(" — %d targeting rules", f.RuleCount))
		}
		sb.WriteString("\n")
		if f.Description != "" {
			sb.WriteString(fmt.Sprintf("  Description: %s\n", f.Description))
		}
		if f.Hypothesis != "" {
			sb.WriteString(fmt.Sprintf("  Hypothesis: %s\n", f.Hypothesis))
		}
		if f.Metrics != "" {
			sb.WriteString(fmt.Sprintf("  Success metrics: %s\n", f.Metrics))
		}
		if f.Status != "" {
			sb.WriteString(fmt.Sprintf("  Status: %s\n", f.Status))
		}
	}

	if len(pc.Members) > 0 {
		sb.WriteString(fmt.Sprintf("\n## Team (%d members)\n", len(pc.Members)))
		for _, m := range pc.Members {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", m.Name, m.Role))
		}
	}

	if len(pc.RecentAudit) > 0 {
		sb.WriteString(fmt.Sprintf("\n## Recent Activity (%d entries)\n", len(pc.RecentAudit)))
		for _, a := range pc.RecentAudit {
			sb.WriteString(fmt.Sprintf("- %s %s %s\n", a.Action, a.EntityType, a.EntityID))
		}
	}

	return sb.String()
}
