package evaluation

import "testing"

func TestOperatorEquals(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"match", "pt-BR", "pt-BR", true, true},
		{"no match", "en-US", "pt-BR", true, false},
		{"not exists", "", "pt-BR", false, false},
		{"empty match", "", "", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpEquals, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("equals(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorNotEquals(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"different", "en-US", "pt-BR", true, true},
		{"same", "pt-BR", "pt-BR", true, false},
		{"not exists returns true", "", "pt-BR", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpNotEquals, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("not_equals(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorContains(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"substring match", "hello world", "world", true, true},
		{"no match", "hello", "world", true, false},
		{"not exists", "", "world", false, false},
		{"empty contains empty", "", "", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpContains, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("contains(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorNotContains(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"no substring", "hello", "world", true, true},
		{"has substring", "hello world", "world", true, false},
		{"not exists returns true", "", "world", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpNotContains, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("not_contains(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorStartsWith(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"match", "pt-BR", "pt", true, true},
		{"no match", "en-US", "pt", true, false},
		{"not exists", "", "pt", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpStartsWith, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("starts_with(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorEndsWith(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"match", "pt-BR", "BR", true, true},
		{"no match", "pt-BR", "US", true, false},
		{"not exists", "", "BR", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpEndsWith, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("ends_with(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorIn(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"in list", "BR", `["BR","US","MX"]`, true, true},
		{"not in list", "FR", `["BR","US","MX"]`, true, false},
		{"not exists", "", `["BR"]`, false, false},
		{"invalid json", "BR", "not json", true, false},
		{"empty list", "BR", `[]`, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpIn, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("in(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorNotIn(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"not in list", "FR", `["BR","US"]`, true, true},
		{"in list", "BR", `["BR","US"]`, true, false},
		{"not exists returns true", "", `["BR"]`, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpNotIn, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("not_in(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorGT(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"greater", "10", "5", true, true},
		{"equal", "5", "5", true, false},
		{"less", "3", "5", true, false},
		{"float greater", "10.5", "10.1", true, true},
		{"not exists", "", "5", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpGT, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("gt(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorGTE(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"greater", "10", "5", true, true},
		{"equal", "5", "5", true, true},
		{"less", "3", "5", true, false},
		{"not exists", "", "5", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpGTE, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("gte(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorLT(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"less", "3", "5", true, true},
		{"equal", "5", "5", true, false},
		{"greater", "10", "5", true, false},
		{"not exists", "", "5", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpLT, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("lt(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorLTE(t *testing.T) {
	tests := []struct {
		name   string
		attr   string
		cond   string
		exists bool
		want   bool
	}{
		{"less", "3", "5", true, true},
		{"equal", "5", "5", true, true},
		{"greater", "10", "5", true, false},
		{"not exists", "", "5", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalOperator(OpLTE, tt.attr, tt.cond, tt.exists); got != tt.want {
				t.Errorf("lte(%q, %q, exists=%v) = %v, want %v", tt.attr, tt.cond, tt.exists, got, tt.want)
			}
		})
	}
}

func TestOperatorExists(t *testing.T) {
	if got := EvalOperator(OpExists, "any", "", true); got != true {
		t.Error("exists with attribute present should return true")
	}
	if got := EvalOperator(OpExists, "", "", false); got != false {
		t.Error("exists with attribute absent should return false")
	}
}

func TestOperatorNotExists(t *testing.T) {
	if got := EvalOperator(OpNotExists, "", "", false); got != true {
		t.Error("not_exists with attribute absent should return true")
	}
	if got := EvalOperator(OpNotExists, "any", "", true); got != false {
		t.Error("not_exists with attribute present should return false")
	}
}

func TestUnknownOperator(t *testing.T) {
	if got := EvalOperator("unknown_op", "a", "b", true); got != false {
		t.Error("unknown operator should return false")
	}
}
