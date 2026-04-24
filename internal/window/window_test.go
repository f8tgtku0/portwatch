package window_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/window"
)

func TestCount_EmptyWindow(t *testing.T) {
	c := window.New(time.Minute, 0)
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCount_WithinWindow(t *testing.T) {
	c := window.New(time.Minute, 0)
	now := time.Now()
	c.RecordAt(now.Add(-10 * time.Second))
	c.RecordAt(now.Add(-5 * time.Second))
	c.RecordAt(now)
	if got := c.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestCount_PrunesOldEvents(t *testing.T) {
	c := window.New(30*time.Second, 0)
	now := time.Now()
	c.RecordAt(now.Add(-60 * time.Second)) // outside window
	c.RecordAt(now.Add(-10 * time.Second)) // inside window
	if got := c.Count(); got != 1 {
		t.Fatalf("expected 1 after pruning, got %d", got)
	}
}

func TestExceeds_BelowThreshold(t *testing.T) {
	c := window.New(time.Minute, 0)
	c.Record()
	c.Record()
	if c.Exceeds(5) {
		t.Fatal("expected Exceeds to return false")
	}
}

func TestExceeds_AboveThreshold(t *testing.T) {
	c := window.New(time.Minute, 0)
	for i := 0; i < 6; i++ {
		c.Record()
	}
	if !c.Exceeds(5) {
		t.Fatal("expected Exceeds to return true")
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	c := window.New(time.Minute, 0)
	c.Record()
	c.Record()
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestMiddleware_Apply_RecordsAndPassesThrough(t *testing.T) {
	var buf bytes.Buffer
	c := window.New(time.Minute, 0)
	mw := window.NewMiddleware(c, 100, &buf)
	changes := []state.Change{
		{Port: 8080, Action: state.Opened},
		{Port: 9090, Action: state.Closed},
	}
	out := mw.Apply(changes)
	if len(out) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out))
	}
	if got := c.Count(); got != 2 {
		t.Fatalf("expected counter to be 2, got %d", got)
	}
}

func TestMiddleware_Apply_LogsWhenThresholdExceeded(t *testing.T) {
	var buf bytes.Buffer
	c := window.New(time.Minute, 0)
	mw := window.NewMiddleware(c, 1, &buf)
	changes := []state.Change{
		{Port: 80, Action: state.Opened},
		{Port: 443, Action: state.Opened},
	}
	mw.Apply(changes)
	if buf.Len() == 0 {
		t.Fatal("expected warning to be written when threshold exceeded")
	}
}

func TestMiddleware_Apply_NilCounter_ReturnsAll(t *testing.T) {
	mw := window.NewMiddleware(nil, 5, nil)
	changes := []state.Change{{Port: 22, Action: state.Opened}}
	out := mw.Apply(changes)
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
}
