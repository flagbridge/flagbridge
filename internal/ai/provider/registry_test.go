package provider

import (
	"sort"
	"strings"
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	fake := NewFake("hi")
	r.Register(fake)

	got, err := r.Get("fake")
	if err != nil {
		t.Fatalf("Get(fake): %v", err)
	}
	if got != fake {
		t.Errorf("Get returned different instance")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing provider, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should name the missing provider: %v", err)
	}
}

func TestRegistry_OverwriteByName(t *testing.T) {
	r := NewRegistry()
	first := &FakeProvider{NameValue: "shared", Deltas: []string{"v1"}}
	second := &FakeProvider{NameValue: "shared", Deltas: []string{"v2"}}
	r.Register(first)
	r.Register(second)

	got, err := r.Get("shared")
	if err != nil {
		t.Fatalf("Get(shared): %v", err)
	}
	if got != second {
		t.Errorf("expected second registration to win")
	}
}

func TestRegistry_Names(t *testing.T) {
	r := NewRegistry()
	r.Register(&FakeProvider{NameValue: "a"})
	r.Register(&FakeProvider{NameValue: "b"})
	r.Register(&FakeProvider{NameValue: "c"})

	names := r.Names()
	sort.Strings(names)
	want := []string{"a", "b", "c"}
	if len(names) != len(want) {
		t.Fatalf("got %d names, want %d", len(names), len(want))
	}
	for i := range names {
		if names[i] != want[i] {
			t.Errorf("names[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}
