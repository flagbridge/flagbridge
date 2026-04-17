package context

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderJSON serialises the context as indented JSON — preferred by OpenAI-family models
// which are trained heavily on JSON-structured input.
func RenderJSON(pc PromptContext) (string, error) {
	buf, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("rendering JSON: %w", err)
	}
	return string(buf), nil
}

// RenderXML serialises the context as XML — Anthropic recommends XML tags for structured
// input (see https://docs.anthropic.com/claude/docs/use-xml-tags). We hand-roll the
// output rather than using encoding/xml with wrapper structs to keep tag names stable
// and attribute ordering predictable for snapshot tests.
func RenderXML(pc PromptContext) string {
	var b strings.Builder
	b.Grow(4096)

	writeLine(&b, 0, fmt.Sprintf(`<context version=%q role=%q>`, pc.Version, pc.Role))

	writeLine(&b, 1, "<project>")
	writeElem(&b, 2, "id", pc.Project.ID)
	writeElem(&b, 2, "name", pc.Project.Name)
	writeElem(&b, 2, "slug", pc.Project.Slug)
	if pc.Project.Description != "" {
		writeElem(&b, 2, "description", pc.Project.Description)
	}
	writeLine(&b, 1, "</project>")

	writeLine(&b, 1, fmt.Sprintf("<flags count=%q>", itoa(len(pc.Flags))))
	for _, f := range pc.Flags {
		writeLine(&b, 2, fmt.Sprintf(`<flag key=%q type=%q>`, f.Key, f.Type))
		writeElem(&b, 3, "name", f.Name)
		if f.Description != "" {
			writeElem(&b, 3, "description", f.Description)
		}
		writeElem(&b, 3, "default_value", f.DefaultValue)
		if len(f.Tags) > 0 {
			writeElem(&b, 3, "tags", strings.Join(f.Tags, ","))
		}
		writeElem(&b, 3, "updated_at", f.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"))
		writeLine(&b, 2, "</flag>")
	}
	writeLine(&b, 1, "</flags>")

	if len(pc.Rules) > 0 {
		writeLine(&b, 1, fmt.Sprintf("<rules count=%q>", itoa(len(pc.Rules))))
		for _, r := range pc.Rules {
			writeLine(&b, 2, fmt.Sprintf(`<rule flag=%q environment=%q priority=%q enabled=%q>`,
				r.FlagKey, r.Environment, itoa(r.Priority), boolStr(r.Enabled)))
			if r.Name != "" {
				writeElem(&b, 3, "name", r.Name)
			}
			writeElem(&b, 3, "conditions", r.Conditions)
			writeElem(&b, 3, "value", r.Value)
			writeLine(&b, 2, "</rule>")
		}
		writeLine(&b, 1, "</rules>")
	}

	if len(pc.ProductCards) > 0 {
		writeLine(&b, 1, fmt.Sprintf("<product_cards count=%q>", itoa(len(pc.ProductCards))))
		for _, c := range pc.ProductCards {
			writeLine(&b, 2, fmt.Sprintf(`<card flag=%q status=%q>`, c.FlagKey, c.Status))
			if c.Hypothesis != "" {
				writeElem(&b, 3, "hypothesis", c.Hypothesis)
			}
			if c.SuccessMetrics != "" {
				writeElem(&b, 3, "success_metrics", c.SuccessMetrics)
			}
			if c.GoNoGo != "" {
				writeElem(&b, 3, "go_no_go", c.GoNoGo)
			}
			writeLine(&b, 2, "</card>")
		}
		writeLine(&b, 1, "</product_cards>")
	}

	if len(pc.RecentAudit) > 0 {
		writeLine(&b, 1, fmt.Sprintf("<recent_audit count=%q>", itoa(len(pc.RecentAudit))))
		for _, a := range pc.RecentAudit {
			attrs := fmt.Sprintf(`action=%q entity_type=%q entity_id=%q at=%q`,
				a.Action, a.EntityType, a.EntityID, a.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"))
			if a.ActorName != "" {
				attrs += fmt.Sprintf(` actor=%q`, a.ActorName)
			}
			writeLine(&b, 2, fmt.Sprintf("<entry %s/>", attrs))
		}
		writeLine(&b, 1, "</recent_audit>")
	}

	if pc.Truncation != nil {
		writeLine(&b, 1, "<truncation>")
		if pc.Truncation.FlagsTruncated {
			writeLine(&b, 2, fmt.Sprintf(`<flags_truncated total=%q/>`, itoa(pc.Truncation.TotalFlags)))
		}
		if pc.Truncation.RulesTruncated {
			writeLine(&b, 2, fmt.Sprintf(`<rules_truncated total=%q/>`, itoa(pc.Truncation.TotalRules)))
		}
		if pc.Truncation.AuditEntriesTruncated {
			writeLine(&b, 2, fmt.Sprintf(`<audit_entries_truncated total=%q/>`, itoa(pc.Truncation.TotalAuditEntries)))
		}
		writeLine(&b, 1, "</truncation>")
	}

	writeLine(&b, 0, "</context>")
	return b.String()
}

func writeLine(b *strings.Builder, depth int, s string) {
	for i := 0; i < depth; i++ {
		b.WriteString("  ")
	}
	b.WriteString(s)
	b.WriteByte('\n')
}

func writeElem(b *strings.Builder, depth int, tag, value string) {
	for i := 0; i < depth; i++ {
		b.WriteString("  ")
	}
	fmt.Fprintf(b, "<%s>%s</%s>\n", tag, escapeXML(value), tag)
}

func escapeXML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return r.Replace(s)
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
