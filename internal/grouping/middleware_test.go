package grouping_test

import (
	"testing"

	"github.com/user/portwatch/internal/grouping"
	"github.com/user/portwatch/internal/state"
)

func makeChanges(ports ...int) []state.Change {
	out := make([]state.Change, len(ports))
	for i, p := range ports {
		out[i] = state.Change{Port: p, Action: state.Opened}
	}
	return out
}

func TestMiddleware_Apply_AnnotatesLabel(t *testing.T) {
	g := grouping.New([]grouping.Group{
		{Name: "web", Ports: []int{443}},
	})
	mw := grouping.NewMiddleware(g)
	out := mw.Apply(makeChanges(443, 8080))
	if len(out) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out))
	}
	if out[0].Label != "web" {
		t.Errorf("expected web label on port 443, got %q", out[0].Label)
	}
	if out[1].Label != "unknown" {
		t.Errorf("expected unknown label on port 8080, got %q", out[1].Label)
	}
}

func TestMiddleware_Apply_NilGrouper_ReturnsAll(t *testing.T) {
	mw := grouping.NewMiddleware(nil)
	changes := makeChanges(80, 443)
	out := mw.Apply(changes)
	if len(out) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out))
	}
}

func TestMiddleware_Apply_EmptyChanges(t *testing.T) {
	g := grouping.New(nil)
	mw := grouping.NewMiddleware(g)
	out := mw.Apply([]state.Change{})
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}
