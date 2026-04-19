package snapshot_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func TestRecorder_Apply_RecordsPorts(t *testing.T) {
	snap := snapshot.New(tempFile(t))
	rec := snapshot.NewRecorder(snap)

	ports := []scanner.Port{{Number: 80}, {Number: 443}}
	out := rec.Apply(ports)

	if len(out) != len(ports) {
		t.Fatalf("Apply changed length: got %d want %d", len(out), len(ports))
	}
	last := snap.Last()
	if last == nil {
		t.Fatal("expected snapshot to be recorded")
	}
	if len(last.Ports) != 2 {
		t.Errorf("snapshot ports = %d, want 2", len(last.Ports))
	}
}

func TestRecorder_Apply_NilSnapshot_ReturnsAll(t *testing.T) {
	rec := snapshot.NewRecorder(nil)
	ports := []scanner.Port{{Number: 22}}
	out := rec.Apply(ports)
	if len(out) != 1 {
		t.Errorf("expected 1 port, got %d", len(out))
	}
}

func TestRecorder_Apply_EmptyPorts(t *testing.T) {
	snap := snapshot.New(tempFile(t))
	rec := snapshot.NewRecorder(snap)
	out := rec.Apply([]scanner.Port{})
	if len(out) != 0 {
		t.Errorf("expected 0 ports, got %d", len(out))
	}
	last := snap.Last()
	if last == nil {
		t.Fatal("expected snapshot entry even for empty ports")
	}
}
