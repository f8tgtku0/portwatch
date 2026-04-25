package acknowledge_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/acknowledge"
)

func TestIsAcknowledged_WithinTTL(t *testing.T) {
	s := acknowledge.New()
	s.Acknowledge(8080, "opened", "planned maintenance", time.Minute)
	if !s.IsAcknowledged(8080, "opened") {
		t.Fatal("expected port 8080 opened to be acknowledged")
	}
}

func TestIsAcknowledged_Expired(t *testing.T) {
	s := acknowledge.New()
	s.Acknowledge(8080, "opened", "", -time.Second)
	if s.IsAcknowledged(8080, "opened") {
		t.Fatal("expected expired acknowledgement to return false")
	}
}

func TestIsAcknowledged_NotPresent(t *testing.T) {
	s := acknowledge.New()
	if s.IsAcknowledged(9090, "closed") {
		t.Fatal("expected unknown port to return false")
	}
}

func TestRevoke_RemovesEntry(t *testing.T) {
	s := acknowledge.New()
	s.Acknowledge(443, "closed", "", time.Hour)
	s.Revoke(443, "closed")
	if s.IsAcknowledged(443, "closed") {
		t.Fatal("expected revoked acknowledgement to return false")
	}
}

func TestPrune_RemovesExpired(t *testing.T) {
	s := acknowledge.New()
	s.Acknowledge(80, "opened", "", -time.Second)  // expired
	s.Acknowledge(443, "opened", "", time.Hour)    // active
	s.Prune()

	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active entry after prune, got %d", len(active))
	}
	if active[0].Port != 443 {
		t.Fatalf("expected port 443, got %d", active[0].Port)
	}
}

func TestActive_ReturnsCopy(t *testing.T) {
	s := acknowledge.New()
	s.Acknowledge(22, "opened", "ssh ok", time.Hour)
	s.Acknowledge(25, "opened", "smtp ok", time.Hour)

	active := s.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active entries, got %d", len(active))
	}

	// Mutating the returned slice must not affect the store.
	active[0].Note = "tampered"
	for _, e := range s.Active() {
		if e.Note == "tampered" {
			t.Fatal("mutating returned slice affected internal store")
		}
	}
}
