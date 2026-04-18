package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	th := throttle.New(time.Second, 2)
	if !th.Allow("port:8080:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_BurstRespected(t *testing.T) {
	th := throttle.New(time.Second, 2)
	key := "port:9090:opened"
	if !th.Allow(key) {
		t.Fatal("expected 1st call allowed")
	}
	if !th.Allow(key) {
		t.Fatal("expected 2nd call allowed")
	}
	if th.Allow(key) {
		t.Fatal("expected 3rd call to be throttled")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(time.Second, 1)
	if !th.Allow("a") {
		t.Fatal("expected a allowed")
	}
	if !th.Allow("b") {
		t.Fatal("expected b allowed independently")
	}
}

func TestAllow_AllowedAfterWindowExpires(t *testing.T) {
	th := throttle.New(50*time.Millisecond, 1)
	key := "port:3000:closed"
	th.Allow(key)
	time.Sleep(60 * time.Millisecond)
	if !th.Allow(key) {
		t.Fatal("expected call allowed after window expired")
	}
}

func TestReset_ClearsState(t *testing.T) {
	th := throttle.New(time.Second, 1)
	th.Allow("x")
	th.Allow("y")
	th.Reset()
	if th.Len() != 0 {
		t.Fatalf("expected 0 keys after reset, got %d", th.Len())
	}
}

func TestLen_TracksKeys(t *testing.T) {
	th := throttle.New(time.Second, 5)
	th.Allow("a")
	th.Allow("b")
	th.Allow("c")
	if th.Len() != 3 {
		t.Fatalf("expected 3 keys, got %d", th.Len())
	}
}
