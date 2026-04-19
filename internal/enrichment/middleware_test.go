package enrichment_test

import (
	"testing"

	"github.com/user/portwatch/internal/enrichment"
	"github.com/user/portwatch/internal/state"
)

func makeChanges(ports ...int) []state.Change {
	out := make([]state.Change, 0, len(ports))
	for _, p := range ports {
		out = append(out, state.Change{Port: p, Action: state.Opened})
	}
	return out
}

func TestMiddleware_EnrichesWellKnownPort(t *testing.T) {
	l := enrichment.New(nil)
	m := enrichment.NewMiddleware(l)

	result := m.Enrich(makeChanges(80))
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Service == "" {
		t.Error("expected non-empty service label for port 80")
	}
}

func TestMiddleware_UnknownPort_EmptyLabel(t *testing.T) {
	l := enrichment.New(nil)
	m := enrichment.NewMiddleware(l)

	result := m.Enrich(makeChanges(39999))
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Service != "" {
		t.Errorf("expected empty label for unknown port, got %q", result[0].Service)
	}
}

func TestMiddleware_NilLookup_ReturnsEmptyLabels(t *testing.T) {
	m := enrichment.NewMiddleware(nil)
	result := m.Enrich(makeChanges(22, 443))
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	for _, r := range result {
		if r.Service != "" {
			t.Errorf("expected empty label with nil lookup, got %q", r.Service)
		}
	}
}

func TestMiddleware_EmptyChanges_ReturnsEmpty(t *testing.T) {
	l := enrichment.New(nil)
	m := enrichment.NewMiddleware(l)
	result := m.Enrich([]state.Change{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestMiddleware_PreservesChangeFields(t *testing.T) {
	l := enrichment.New(nil)
	m := enrichment.NewMiddleware(l)
	changes := []state.Change{{Port: 22, Action: state.Opened}}
	result := m.Enrich(changes)
	if result[0].Port != 22 {
		t.Errorf("expected port 22, got %d", result[0].Port)
	}
	if result[0].Action != state.Opened {
		t.Errorf("expected Opened action")
	}
}
