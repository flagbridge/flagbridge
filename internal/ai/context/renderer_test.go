package context

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var updateGolden = flag.Bool("update", false, "update golden files with current renderer output")

func fixedTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parsing time %q: %v", s, err)
	}
	return v
}

func sampleContext(t *testing.T) PromptContext {
	t.Helper()
	return PromptContext{
		Version: CurrentVersion,
		Role:    RoleProduct,
		Project: ProjectContext{
			ID:          "11111111-1111-4111-8111-111111111111",
			Name:        "Checkout",
			Slug:        "checkout",
			Description: "Checkout funnel experiments",
		},
		Flags: []FlagContext{
			{
				Key:          "dark-mode",
				Name:         "Dark mode",
				Description:  "Enable dark theme in admin",
				Type:         "boolean",
				DefaultValue: "false",
				Tags:         []string{"ui", "theme"},
				UpdatedAt:    fixedTime(t, "2026-04-14T10:00:00Z"),
			},
			{
				Key:          "checkout-v2",
				Name:         "Checkout v2",
				Type:         "string",
				DefaultValue: `"legacy"`,
				UpdatedAt:    fixedTime(t, "2026-04-15T12:00:00Z"),
			},
		},
		Rules: []RuleContext{
			{
				FlagKey:     "checkout-v2",
				Environment: "production",
				Name:        "Beta cohort",
				Priority:    10,
				Enabled:     true,
				Conditions:  `[{"attr":"user.plan","op":"eq","val":"beta"}]`,
				Value:       `"v2"`,
			},
		},
		ProductCards: []ProductCardContext{
			{
				FlagKey:        "checkout-v2",
				Hypothesis:     "New funnel lifts conversion by 8%",
				SuccessMetrics: "conversion_rate, AOV",
				GoNoGo:         "Ship if conversion_rate > 1.05x baseline at p95",
				Status:         "in_progress",
			},
		},
		RecentAudit: []AuditEntry{
			{
				Action:     "flag.toggle",
				EntityType: "flag",
				EntityID:   "22222222-2222-4222-8222-222222222222",
				ActorName:  "Gabriel Gripp",
				CreatedAt:  fixedTime(t, "2026-04-16T09:30:00Z"),
			},
		},
	}
}

func TestRenderJSON_Golden(t *testing.T) {
	pc := sampleContext(t)
	got, err := RenderJSON(pc)
	if err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	compareGolden(t, "sample-context.json.golden", got)
}

func TestRenderXML_Golden(t *testing.T) {
	pc := sampleContext(t)
	got := RenderXML(pc)
	compareGolden(t, "sample-context.xml.golden", got)
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	pc := sampleContext(t)
	got, err := RenderJSON(pc)
	if err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	var back PromptContext
	if err := json.Unmarshal([]byte(got), &back); err != nil {
		t.Fatalf("unmarshal roundtrip: %v", err)
	}
	if back.Version != CurrentVersion {
		t.Errorf("version roundtrip: got %q want %q", back.Version, CurrentVersion)
	}
	if len(back.Flags) != len(pc.Flags) {
		t.Errorf("flags roundtrip: got %d want %d", len(back.Flags), len(pc.Flags))
	}
}

func TestRenderXML_EscapesSpecialChars(t *testing.T) {
	pc := PromptContext{
		Version: CurrentVersion,
		Role:    RoleEngineer,
		Project: ProjectContext{ID: "p1", Name: "A & B <test>", Slug: "ab"},
		Flags: []FlagContext{
			{Key: "k", Name: "<script>", Type: "boolean", DefaultValue: "false", UpdatedAt: fixedTime(t, "2026-04-17T00:00:00Z")},
		},
	}
	got := RenderXML(pc)
	if strings.Contains(got, "<script>") {
		t.Errorf("unescaped < > in XML output: %s", got)
	}
	if !strings.Contains(got, "&amp;") {
		t.Errorf("expected ampersand escape: %s", got)
	}
}

func TestRenderXML_EngineerRoleOmitsProductCards(t *testing.T) {
	// Even if the caller mistakenly includes cards, the render shouldn't crash —
	// but a correct builder won't populate them for engineer role. Here we verify
	// that an engineer context with empty ProductCards emits no <product_cards> tag.
	pc := sampleContext(t)
	pc.Role = RoleEngineer
	pc.ProductCards = nil
	got := RenderXML(pc)
	if strings.Contains(got, "<product_cards") {
		t.Errorf("engineer role should omit <product_cards> tag: %s", got)
	}
}

func TestRenderXML_TruncationSurfacedInOutput(t *testing.T) {
	pc := sampleContext(t)
	pc.Truncation = &TruncationInfo{
		FlagsTruncated: true,
		TotalFlags:     137,
	}
	got := RenderXML(pc)
	if !strings.Contains(got, "<truncation>") {
		t.Errorf("missing <truncation> block: %s", got)
	}
	if !strings.Contains(got, `total="137"`) {
		t.Errorf("missing total attr: %s", got)
	}
}

func compareGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *updateGolden {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("writing golden %s: %v", name, err)
		}
		t.Logf("updated golden %s", name)
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading golden %s: %v (run with -update to generate)", name, err)
	}
	if got != string(want) {
		t.Errorf("golden %s mismatch.\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}
