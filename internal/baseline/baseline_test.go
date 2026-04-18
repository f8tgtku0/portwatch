package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestAdd_Contains(t *testing.T) {
	b := baseline.New(tempFile(t))
	b.Add(8080)
	if !b.Contains(8080) {
		t.Fatal("expected 8080 to be in baseline")
	}
}

func TestRemove(t *testing.T) {
	b := baseline.New(tempFile(t))
	b.Add(443)
	b.Remove(443)
	if b.Contains(443) {
		t.Fatal("expected 443 to be removed")
	}
}

func TestSave_AndLoad(t *testing.T) {
	path := tempFile(t)
	b := baseline.New(path)
	b.Add(22)
	b.Add(80)
	if err := b.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	b2 := baseline.New(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if !b2.Contains(22) || !b2.Contains(80) {
		t.Fatal("expected loaded baseline to contain saved ports")
	}
}

func TestLoad_NoFile(t *testing.T) {
	b := baseline.New("/nonexistent/path/baseline.json")
	if err := b.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestPorts_ReturnsCopy(t *testing.T) {
	b := baseline.New(tempFile(t))
	b.Add(9000)
	ports := b.Ports()
	if len(ports) != 1 || ports[0] != 9000 {
		t.Fatalf("unexpected ports: %v", ports)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tempFile(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	b := baseline.New(path)
	if err := b.Load(); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
