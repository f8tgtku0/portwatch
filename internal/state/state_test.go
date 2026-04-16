package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "state.json")
}

func TestSave_AndLoad(t *testing.T) {
	path := tempFile(t)
	s := state.New(path)

	ports := []int{80, 443, 8080}
	if err := s.Save(ports); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if len(snap.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	for i, p := range ports {
		if snap.Ports[i] != p {
			t.Errorf("port[%d]: expected %d, got %d", i, p, snap.Ports[i])
		}
	}
}

func TestLoad_NoFile(t *testing.T) {
	path := tempFile(t)
	s := state.New(path)

	snap, err := s.Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if snap != nil {
		t.Errorf("expected nil snapshot when no file exists, got %+v", snap)
	}
}

func TestLoad_PersistsAcrossInstances(t *testing.T) {
	path := tempFile(t)

	s1 := state.New(path)
	if err := s1.Save([]int{22, 3306}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	s2 := state.New(path)
	snap, err := s2.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if snap == nil || len(snap.Ports) != 2 {
		t.Errorf("expected 2 ports from new store instance, got %v", snap)
	}
}

func TestClear_RemovesFile(t *testing.T) {
	path := tempFile(t)
	s := state.New(path)

	_ = s.Save([]int{9000})
	if err := s.Clear(); err != nil {
		t.Fatalf("Clear() error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected state file to be removed after Clear()")
	}

	snap, err := s.Load()
	if err != nil {
		t.Fatalf("Load() after Clear() error: %v", err)
	}
	if snap != nil {
		t.Errorf("expected nil snapshot after Clear(), got %+v", snap)
	}
}

func TestClear_NoFile_NoError(t *testing.T) {
	path := tempFile(t)
	s := state.New(path)

	if err := s.Clear(); err != nil {
		t.Errorf("Clear() on non-existent file should not error, got: %v", err)
	}
}

func TestSave_EmptyPorts(t *testing.T) {
	path := tempFile(t)
	s := state.New(path)

	if err := s.Save([]int{}); err != nil {
		t.Fatalf("Save() with empty ports error: %v", err)
	}

	snap, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected 0 ports, got %d", len(snap.Ports))
	}
}
