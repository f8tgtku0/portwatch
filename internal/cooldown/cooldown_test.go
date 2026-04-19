package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
	"github.com/user/portwatch/internal/state"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	if !c.Allow(8080, state.Opened) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinWindowBlocked(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow(8080, state.Opened)
	if c.Allow(8080, state.Opened) {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_AllowedAfterWindowExpires(t *testing.T) {
	c := cooldown.New(20 * time.Millisecond)
	c.Allow(8080, state.Opened)
	time.Sleep(30 * time.Millisecond)
	if !c.Allow(8080, state.Opened) {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestAllow_DifferentActionsAreIndependent(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow(8080, state.Opened)
	if !c.Allow(8080, state.Closed) {
		t.Fatal("expected different action on same port to be allowed")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow(8080, state.Opened)
	if !c.Allow(9090, state.Opened) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_ClearsState(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow(8080, state.Opened)
	c.Reset()
	if !c.Allow(8080, state.Opened) {
		t.Fatal("expected allow after reset")
	}
}

func TestReset_LenIsZeroAfterReset(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow(8080, state.Opened)
	c.Allow(9090, state.Closed)
	c.Reset()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", c.Len())
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	c := cooldown.New(20 * time.Millisecond)
	c.Allow(8080, state.Opened)
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", c.Len())
	}
	time.Sleep(30 * time.Millisecond)
	c.Prune()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after prune, got %d", c.Len())
	}
}
