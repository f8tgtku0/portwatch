package redact_test

import (
	"testing"

	"github.com/user/portwatch/internal/redact"
	"github.com/user/portwatch/internal/state"
)

func changes(ports ...int) []state.Change {
	out := make([]state.Change, len(ports))
	for i, p := range ports {
		out[i] = state.Change{Port: p, Action: state.Opened}
	}
	return out
}

func TestApply_HidesMaskedPorts(t *testing.T) {
	r := redact.New([]int{22, 3306})
	visible, count := r.Apply(changes(80, 22, 443, 3306))
	if count != 2 {
		t.Fatalf("expected 2 redacted, got %d", count)
	}
	if len(visible) != 2 {
		t.Fatalf("expected 2 visible, got %d", len(visible))
	}
	for _, c := range visible {
		if c.Port == 22 || c.Port == 3306 {
			t.Errorf("masked port %d leaked into visible set", c.Port)
		}
	}
}

func TestApply_NoMaskedPorts(t *testing.T) {
	r := redact.New(nil)
	visible, count := r.Apply(changes(80, 443))
	if count != 0 {
		t.Fatalf("expected 0 redacted, got %d", count)
	}
	if len(visible) != 2 {
		t.Fatalf("expected 2 visible, got %d", len(visible))
	}
}

func TestAdd_MasksNewPort(t *testing.T) {
	r := redact.New(nil)
	r.Add(8080)
	_, count := r.Apply(changes(8080))
	if count != 1 {
		t.Fatalf("expected 1 redacted after Add, got %d", count)
	}
}

func TestRemove_UnmasksPort(t *testing.T) {
	r := redact.New([]int{9200})
	r.Remove(9200)
	visible, count := r.Apply(changes(9200))
	if count != 0 {
		t.Fatalf("expected 0 redacted after Remove, got %d", count)
	}
	if len(visible) != 1 {
		t.Fatalf("expected 1 visible after Remove, got %d", len(visible))
	}
}

func TestIsMasked(t *testing.T) {
	r := redact.New([]int{5432})
	if !r.IsMasked(5432) {
		t.Error("expected 5432 to be masked")
	}
	if r.IsMasked(80) {
		t.Error("expected 80 to not be masked")
	}
}

func TestSummary_ZeroCount(t *testing.T) {
	if s := redact.Summary(0); s != "" {
		t.Errorf("expected empty string for zero count, got %q", s)
	}
}

func TestSummary_NonZero(t *testing.T) {
	s := redact.Summary(3)
	if s == "" {
		t.Error("expected non-empty summary for count=3")
	}
}
