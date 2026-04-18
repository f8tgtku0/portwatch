package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

var (
	now   = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	early = now.Add(-1 * time.Hour)
	late  = now.Add(1 * time.Hour)
)

func TestIsSuppressed_WithinWindow(t *testing.T) {
	s := suppress.New()
	s.Add(early, late, "maintenance")
	ok, w := s.IsSuppressed(now)
	if !ok {
		t.Fatal("expected suppressed")
	}
	if w.Reason != "maintenance" {
		t.Errorf("unexpected reason: %s", w.Reason)
	}
}

func TestIsSuppressed_OutsideWindow(t *testing.T) {
	s := suppress.New()
	s.Add(late, late.Add(time.Hour), "later")
	ok, _ := s.IsSuppressed(now)
	if ok {
		t.Fatal("expected not suppressed")
	}
}

func TestIsSuppressed_NoWindows(t *testing.T) {
	s := suppress.New()
	ok, _ := s.IsSuppressed(now)
	if ok {
		t.Fatal("expected not suppressed")
	}
}

func TestPrune_RemovesExpired(t *testing.T) {
	s := suppress.New()
	s.Add(early, now.Add(-30*time.Minute), "expired")
	s.Add(early, late, "active")
	s.Prune(now)
	windows := s.Active()
	if len(windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(windows))
	}
	if windows[0].Reason != "active" {
		t.Errorf("wrong window retained: %s", windows[0].Reason)
	}
}

func TestActive_ReturnsCopy(t *testing.T) {
	s := suppress.New()
	s.Add(early, late, "w1")
	a := s.Active()
	a[0].Reason = "mutated"
	b := s.Active()
	if b[0].Reason == "mutated" {
		t.Error("Active() should return a copy")
	}
}
