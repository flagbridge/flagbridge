package context

import "testing"

func TestIncludesProductCards(t *testing.T) {
	cases := []struct {
		role string
		want bool
	}{
		{RoleProduct, true},
		{RoleAdmin, true},
		{RoleEngineer, false},
		{RoleViewer, false},
		{"", false},
		{"unknown", false},
	}
	for _, tc := range cases {
		t.Run(tc.role, func(t *testing.T) {
			if got := IncludesProductCards(tc.role); got != tc.want {
				t.Errorf("IncludesProductCards(%q) = %v, want %v", tc.role, got, tc.want)
			}
		})
	}
}

func TestBudgetConstants(t *testing.T) {
	// Guard against accidental relaxation. The budgets come from ADR-0002;
	// changing them requires an ADR update.
	if MaxFlags != 50 {
		t.Errorf("MaxFlags changed without ADR update: got %d want 50", MaxFlags)
	}
	if MaxAuditEntries != 20 {
		t.Errorf("MaxAuditEntries changed without ADR update: got %d want 20", MaxAuditEntries)
	}
	if MaxRules != 100 {
		t.Errorf("MaxRules changed without ADR update: got %d want 100", MaxRules)
	}
}
