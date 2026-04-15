package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestRecord_AddsEntry(t *testing.T) {
	h, err := New(tempFile(t), 100)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := h.Record(8080, "opened"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Event != "opened" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestRecord_RotatesAtMaxSize(t *testing.T) {
	h, err := New(tempFile(t), 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := h.Record(8000+i, "opened"); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	entries := h.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries after rotation, got %d", len(entries))
	}
	if entries[0].Port != 8002 {
		t.Errorf("expected oldest retained port 8002, got %d", entries[0].Port)
	}
}

func TestRecord_PersistsToDisk(t *testing.T) {
	path := tempFile(t)
	h, err := New(path, 100)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = h.Record(9090, "closed")

	h2, err := New(path, 100)
	if err != nil {
		t.Fatalf("New (reload): %v", err)
	}
	entries := h2.Entries()
	if len(entries) != 1 || entries[0].Port != 9090 {
		t.Errorf("persisted entry mismatch: %+v", entries)
	}
}

func TestClear_RemovesEntries(t *testing.T) {
	path := tempFile(t)
	h, err := New(path, 100)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = h.Record(1234, "opened")
	if err := h.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if len(h.Entries()) != 0 {
		t.Error("expected no entries after Clear")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
}

func TestNew_MissingFileIsOk(t *testing.T) {
	h, err := New(tempFile(t), 50)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries()) != 0 {
		t.Error("expected empty history for new file")
	}
}
