package pause_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/pause"
)

func TestIsPaused_InitiallyFalse(t *testing.T) {
	p := pause.New()
	if p.IsPaused() {
		t.Fatal("expected not paused initially")
	}
}

func TestPause_SetsState(t *testing.T) {
	p := pause.New()
	p.Pause(5 * time.Minute)
	if !p.IsPaused() {
		t.Fatal("expected paused after Pause()")
	}
}

func TestResume_ClearsPause(t *testing.T) {
	p := pause.New()
	p.Pause(5 * time.Minute)
	p.Resume()
	if p.IsPaused() {
		t.Fatal("expected not paused after Resume()")
	}
}

func TestPause_ExpiresAfterDuration(t *testing.T) {
	p := pause.New()
	// Use a very short duration and advance the clock via the hook.
	now := time.Now()
	p.SetNow(func() time.Time { return now })
	p.Pause(100 * time.Millisecond)

	// Still paused at t=0.
	if !p.IsPaused() {
		t.Fatal("expected paused immediately after Pause()")
	}

	// Advance clock past the deadline.
	p.SetNow(func() time.Time { return now.Add(200 * time.Millisecond) })
	if p.IsPaused() {
		t.Fatal("expected pause to have expired")
	}
}

func TestPause_ExtendDeadline(t *testing.T) {
	p := pause.New()
	now := time.Now()
	p.SetNow(func() time.Time { return now })
	p.Pause(1 * time.Minute)
	p.Pause(10 * time.Minute) // extend

	got := p.Until()
	want := now.Add(10 * time.Minute)
	if !got.Equal(want) {
		t.Fatalf("Until() = %v, want %v", got, want)
	}
}

func TestUntil_ZeroWhenNotPaused(t *testing.T) {
	p := pause.New()
	if !p.Until().IsZero() {
		t.Fatal("expected zero Until when not paused")
	}
}
