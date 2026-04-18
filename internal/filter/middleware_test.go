package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/state"
)

func makeChanges(ports ...int) []state.Change {
	out := make([]state.Change, len(ports))
	for i, p := range ports {
		out[i] = state.Change{Port: p, Opened: true}
	}
	return out
}

func TestMiddleware_Apply_FiltersMatchedPorts(t *testing.T) {
	f, err := filter.New([]string{"8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := filter.NewMiddleware(f)
	result := m.Apply(makeChanges(80, 8080, 443))
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
	for _, c := range result {
		if c.Port == 8080 {
			t.Error("port 8080 should have been filtered out")
		}
	}
}

func TestMiddleware_Apply_NilFilter_ReturnsAll(t *testing.T) {
	m := filter.NewMiddleware(nil)
	changes := makeChanges(80, 443)
	result := m.Apply(changes)
	if len(result) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(result))
	}
}

func TestMiddleware_Apply_RangeFilter(t *testing.T) {
	f, err := filter.New([]string{"8000-9000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := filter.NewMiddleware(f)
	result := m.Apply(makeChanges(80, 8080, 8443, 443))
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestMiddleware_Apply_EmptyChanges(t *testing.T) {
	f, err := filter.New([]string{"80"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := filter.NewMiddleware(f)
	result := m.Apply([]state.Change{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}
