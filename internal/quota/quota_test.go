package quota_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/quota"
	"github.com/user/portwatch/internal/state"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	q := quota.New(3, time.Minute)
	if !q.Allow("slack") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_ExhaustsQuota(t *testing.T) {
	q := quota.New(2, time.Minute)
	if !q.Allow("ch") {
		t.Fatal("call 1 should be allowed")
	}
	if !q.Allow("ch") {
		t.Fatal("call 2 should be allowed")
	}
	if q.Allow("ch") {
		t.Fatal("call 3 should be blocked")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now := time.Now()
	q := quota.New(1, 50*time.Millisecond)
	// inject controllable clock
	q2 := quota.New(1, 50*time.Millisecond)
	_ = now
	_ = q2

	if !q.Allow("x") {
		t.Fatal("first call should be allowed")
	}
	if q.Allow("x") {
		t.Fatal("second call within window should be blocked")
	}
	time.Sleep(60 * time.Millisecond)
	if !q.Allow("x") {
		t.Fatal("call after window reset should be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	q := quota.New(1, time.Minute)
	q.Allow("a")
	if !q.Allow("b") {
		t.Fatal("key b should still have quota")
	}
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	q := quota.New(3, time.Minute)
	if got := q.Remaining("k"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
	q.Allow("k")
	if got := q.Remaining("k"); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}
}

func TestReset_ClearsAllKeys(t *testing.T) {
	q := quota.New(1, time.Minute)
	q.Allow("a")
	q.Allow("b")
	q.Reset()
	if got := q.Remaining("a"); got != 1 {
		t.Fatalf("expected quota restored after reset, got %d", got)
	}
}

func TestMiddleware_Apply_DropsExhaustedChannel(t *testing.T) {
	q := quota.New(1, time.Minute)
	keyer := func(c state.Change) string { return "ch" }
	m := quota.NewMiddleware(q, keyer)

	changes := []state.Change{
		{Port: 80, Action: state.Opened},
		{Port: 443, Action: state.Opened},
	}
	out := m.Apply(changes)
	if len(out) != 1 {
		t.Fatalf("expected 1 change after quota, got %d", len(out))
	}
}

func TestMiddleware_Apply_NilQuota_ReturnsAll(t *testing.T) {
	m := quota.NewMiddleware(nil, func(c state.Change) string { return "" })
	changes := []state.Change{{Port: 22, Action: state.Opened}}
	if got := m.Apply(changes); len(got) != 1 {
		t.Fatal("nil quota should pass all changes through")
	}
}
