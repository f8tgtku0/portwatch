package truncate_test

import (
	"bytes"
	"strings"
	"testing"

	"portwatch/internal/state"
	"portwatch/internal/truncate"
)

func makeChanges(ports ...int) []state.Change {
	out := make([]state.Change, len(ports))
	for i, p := range ports {
		out[i] = state.Change{Port: p, Action: state.Opened}
	}
	return out
}

func TestApply_BelowLimit_ReturnsAll(t *testing.T) {
	tr := truncate.New(5, nil)
	changes := makeChanges(80, 443, 8080)
	got := tr.Apply(changes)
	if len(got) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(got))
	}
}

func TestApply_AtLimit_ReturnsAll(t *testing.T) {
	tr := truncate.New(3, nil)
	changes := makeChanges(80, 443, 8080)
	got := tr.Apply(changes)
	if len(got) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(got))
	}
}

func TestApply_ExceedsLimit_Truncates(t *testing.T) {
	var buf bytes.Buffer
	tr := truncate.New(2, &buf)
	changes := makeChanges(80, 443, 8080, 9090)
	got := tr.Apply(changes)
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
	if got[0].Port != 80 || got[1].Port != 443 {
		t.Errorf("unexpected ports in truncated slice: %v", got)
	}
}

func TestApply_ExceedsLimit_WritesWarning(t *testing.T) {
	var buf bytes.Buffer
	tr := truncate.New(2, &buf)
	tr.Apply(makeChanges(1, 2, 3, 4, 5))
	if !strings.Contains(buf.String(), "truncated") {
		t.Errorf("expected warning message, got: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "3 dropped") {
		t.Errorf("expected dropped count in warning, got: %q", buf.String())
	}
}

func TestNew_ClampsBelowOne(t *testing.T) {
	tr := truncate.New(0, nil)
	if tr.MaxSize() != 1 {
		t.Errorf("expected maxSize=1 for zero input, got %d", tr.MaxSize())
	}
}

func TestNew_NilWriter_DefaultsToStderr(t *testing.T) {
	// Should not panic when writer is nil (defaults to os.Stderr internally).
	tr := truncate.New(1, nil)
	changes := makeChanges(80, 443)
	got := tr.Apply(changes)
	if len(got) != 1 {
		t.Errorf("expected 1 change, got %d", len(got))
	}
}

func TestApply_EmptyChanges_ReturnsEmpty(t *testing.T) {
	tr := truncate.New(10, nil)
	got := tr.Apply(nil)
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
