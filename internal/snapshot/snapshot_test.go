package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestRecord_AndLast(t *testing.T) {
	s := snapshot.New(tempFile(t))
	ports := []scanner.Port{{Number: 80}, {Number: 443}}
	if err := s.Record(ports); err != nil {
		t.Fatalf("Record: %v", err)
	}
	got := s.Last()
	if got == nil {
		t.Fatal("expected non-nil Last")
	}
	if len(got.Ports) != 2 {
		t.Errorf("ports len = %d, want 2", len(got.Ports))
	}
	if got.TakenAt.IsZero() {
		t.Error("TakenAt should not be zero")
	}
}

func TestLoad_NoFile(t *testing.T) {
	s := snapshot.New(tempFile(t))
	if err := s.Load(); err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if s.Last() != nil {
		t.Error("expected nil Last when no file")
	}
}

func TestLoad_PersistsAcrossInstances(t *testing.T) {
	path := tempFile(t)
	s1 := snapshot.New(path)
	ports := []scanner.Port{{Number: 22}}
	_ = s1.Record(ports)

	s2 := snapshot.New(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	got := s2.Last()
	if got == nil || len(got.Ports) != 1 || got.Ports[0].Number != 22 {
		t.Errorf("unexpected Last: %+v", got)
	}
}

func TestClear_RemovesSnapshot(t *testing.T) {
	path := tempFile(t)
	s := snapshot.New(path)
	_ = s.Record([]scanner.Port{{Number: 8080}})
	if err := s.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if s.Last() != nil {
		t.Error("expected nil after Clear")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestRecord_UpdatesTimestamp(t *testing.T) {
	s := snapshot.New(tempFile(t))
	before := time.Now()
	_ = s.Record([]scanner.Port{{Number: 3306}})
	after := time.Now()
	got := s.Last()
	if got.TakenAt.Before(before) || got.TakenAt.After(after) {
		t.Errorf("TakenAt %v not in [%v, %v]", got.TakenAt, before, after)
	}
}
