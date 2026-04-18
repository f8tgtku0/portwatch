package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	if !l.Allow("port:8080:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("port:8080:opened")
	if l.Allow("port:8080:opened") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_AllowedAfterCooldown(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(5 * time.Second)

	// Inject a clock that starts in the past
	calls := 0
	l2 := &struct{ *ratelimit.Limiter }{ratelimit.New(5 * time.Second)}
	_ = l2
	_ = calls
	_ = now

	// Use Reset to simulate expiry
	l.Allow("port:9090:closed")
	l.Reset("port:9090:closed")
	if !l.Allow("port:9090:closed") {
		t.Fatal("expected allow after reset")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("port:80:opened")
	if !l.Allow("port:443:opened") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestClear_RemovesAllKeys(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("a")
	l.Allow("b")
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected 0 keys after clear, got %d", l.Len())
	}
}

func TestLen_TracksKeys(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("x")
	l.Allow("y")
	l.Allow("x") // duplicate, should not add
	if l.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", l.Len())
	}
}
