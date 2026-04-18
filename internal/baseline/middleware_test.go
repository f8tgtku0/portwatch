package baseline_test

import (
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/state"
)

func makeChanges(ports ...int) []state.Change {
	out := make([]state.Change, len(ports))
	for i, p := range ports {
		out[i] = state.Change{Port: p, Action: state.Opened}
	}
	return out
}

func TestMiddleware_FiltersBaselinePorts(t *testing.T) {
	b := baseline.New(tempFile(t))
	b.Add(22)
	b.Add(80)
	mw := baseline.NewMiddleware(b)

	changes := makeChanges(22, 80, 8080)
	result := mw.Apply(changes)
	if len(result) != 1 || result[0].Port != 8080 {
		t.Fatalf("expected only port 8080, got %v", result)
	}
}

func TestMiddleware_NilBaseline_ReturnsAll(t *testing.T) {
	mw := baseline.NewMiddleware(nil)
	changes := makeChanges(22, 443)
	result := mw.Apply(changes)
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestMiddleware_EmptyBaseline_ReturnsAll(t *testing.T) {
	b := baseline.New(tempFile(t))
	mw := baseline.NewMiddleware(b)
	changes := makeChanges(8080, 9090)
	result := mw.Apply(changes)
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestMiddleware_AllFiltered(t *testing.T) {
	b := baseline.New(tempFile(t))
	b.Add(22)
	mw := baseline.NewMiddleware(b)
	result := mw.Apply(makeChanges(22))
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}
