package correlation_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlation"
	"github.com/user/portwatch/internal/state"
)

func TestMiddleware_Apply_NilCorrelator_PassesThrough(t *testing.T) {
	mw := correlation.NewMiddleware(nil)
	in := []state.Change{opened(80), closed(443)}
	out := mw.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestMiddleware_Apply_BuffersWithinWindow(t *testing.T) {
	c := correlation.New(200 * time.Millisecond)
	mw := correlation.NewMiddleware(c)

	out := mw.Apply([]state.Change{opened(80)})
	if len(out) != 0 {
		t.Fatalf("expected 0 (buffered), got %d", len(out))
	}
}

func TestMiddleware_Apply_ReturnsAfterFlush(t *testing.T) {
	c := correlation.New(50 * time.Millisecond)
	mw := correlation.NewMiddleware(c)

	mw.Apply([]state.Change{opened(8080)})
	time.Sleep(80 * time.Millisecond)

	// A second apply with no new changes should drain the ready event.
	out := mw.Apply(nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 drained change, got %d", len(out))
	}
}

func TestMiddleware_Apply_EmptyChanges_ReturnsNil(t *testing.T) {
	c := correlation.New(50 * time.Millisecond)
	mw := correlation.NewMiddleware(c)
	out := mw.Apply(nil)
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}
