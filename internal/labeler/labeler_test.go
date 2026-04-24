package labeler_test

import (
	"testing"

	"github.com/user/portwatch/internal/labeler"
	"github.com/user/portwatch/internal/state"
)

func TestLookup_WellKnownPort(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Lookup(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Lookup(9999); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestLookup_OverrideTakesPriority(t *testing.T) {
	l := labeler.New(map[int]string{80: "custom-http"})
	if got := l.Lookup(80); got != "custom-http" {
		t.Fatalf("expected custom-http, got %q", got)
	}
}

func TestLookup_CustomNewPort(t *testing.T) {
	l := labeler.New(map[int]string{8888: "jupyter"})
	if got := l.Lookup(8888); got != "jupyter" {
		t.Fatalf("expected jupyter, got %q", got)
	}
}

func TestAnnotate_SetsLabel(t *testing.T) {
	l := labeler.New(nil)
	changes := []state.Change{
		{Port: 443, Action: state.Opened},
		{Port: 9999, Action: state.Opened},
	}
	out := l.Annotate(changes)
	if out[0].Label != "https" {
		t.Errorf("port 443: expected https, got %q", out[0].Label)
	}
	if out[1].Label != "" {
		t.Errorf("port 9999: expected empty, got %q", out[1].Label)
	}
}

func TestAnnotate_NilLabeler_ReturnsOriginal(t *testing.T) {
	var l *labeler.Labeler
	changes := []state.Change{{Port: 22, Action: state.Opened}}
	out := l.Annotate(changes)
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
}

func TestAnnotate_EmptyChanges(t *testing.T) {
	l := labeler.New(nil)
	out := l.Annotate([]state.Change{})
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(out))
	}
}
