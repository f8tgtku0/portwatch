package mute_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/mute"
)

func TestIsMuted_ActiveRule(t *testing.T) {
	m := mute.New()
	m.Add(8080, "maintenance", 1*time.Hour)
	if !m.IsMuted(8080) {
		t.Fatal("expected port 8080 to be muted")
	}
}

func TestIsMuted_UnknownPort(t *testing.T) {
	m := mute.New()
	if m.IsMuted(9999) {
		t.Fatal("expected unmuted port to return false")
	}
}

func TestIsMuted_ExpiredRule(t *testing.T) {
	m := mute.New()
	m.Add(443, "expired", -1*time.Second) // already expired
	if m.IsMuted(443) {
		t.Fatal("expected expired rule to be treated as inactive")
	}
}

func TestRemove_ClearsMuteRule(t *testing.T) {
	m := mute.New()
	m.Add(22, "test", 1*time.Hour)
	m.Remove(22)
	if m.IsMuted(22) {
		t.Fatal("expected port 22 to be unmuted after Remove")
	}
}

func TestActive_ReturnsOnlyLiveEntries(t *testing.T) {
	m := mute.New()
	m.Add(80, "live", 1*time.Hour)
	m.Add(81, "expired", -1*time.Second)

	actives := m.Active()
	if len(actives) != 1 {
		t.Fatalf("expected 1 active entry, got %d", len(actives))
	}
	if actives[0].Port != 80 {
		t.Errorf("expected port 80, got %d", actives[0].Port)
	}
}

func TestActive_ReturnsReason(t *testing.T) {
	m := mute.New()
	m.Add(3306, "db maintenance", 30*time.Minute)

	actives := m.Active()
	if len(actives) == 0 {
		t.Fatal("expected at least one active entry")
	}
	if actives[0].Reason != "db maintenance" {
		t.Errorf("unexpected reason: %q", actives[0].Reason)
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	m := mute.New()
	m.Add(8443, "prune-me", -1*time.Second)
	m.Add(8444, "keep-me", 1*time.Hour)
	m.Prune()

	actives := m.Active()
	for _, e := range actives {
		if e.Port == 8443 {
			t.Error("expected pruned port 8443 to be absent")
		}
	}
}

func TestAdd_ReplacesExistingRule(t *testing.T) {
	m := mute.New()
	m.Add(5432, "first", 1*time.Minute)
	m.Add(5432, "second", 2*time.Hour)

	actives := m.Active()
	if len(actives) != 1 {
		t.Fatalf("expected 1 entry after replacement, got %d", len(actives))
	}
	if actives[0].Reason != "second" {
		t.Errorf("expected updated reason 'second', got %q", actives[0].Reason)
	}
}
