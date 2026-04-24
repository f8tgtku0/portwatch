package backoff_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func noSleep(_ time.Duration) {}

func TestDefaultPolicy_Fields(t *testing.T) {
	p := backoff.DefaultPolicy()
	if p.InitialInterval <= 0 {
		t.Fatal("expected positive InitialInterval")
	}
	if p.MaxInterval < p.InitialInterval {
		t.Fatal("MaxInterval must be >= InitialInterval")
	}
	if p.Multiplier <= 1.0 {
		t.Fatal("Multiplier must be > 1.0")
	}
}

func TestBackoff_AttemptZero_ReturnsInitial(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0.0, // no jitter for deterministic test
	}
	got := p.Backoff(0)
	if got != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", got)
	}
}

func TestBackoff_GrowsExponentially(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0.0,
	}
	prev := p.Backoff(0)
	for i := 1; i <= 4; i++ {
		curr := p.Backoff(i)
		if curr <= prev {
			t.Fatalf("attempt %d: expected growth, got %v <= %v", i, curr, prev)
		}
		prev = curr
	}
}

func TestBackoff_CapsAtMaxInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     2 * time.Second,
		Multiplier:      10.0,
		JitterFraction:  0.0,
	}
	got := p.Backoff(10)
	if got > 2*time.Second {
		t.Fatalf("expected <= 2s, got %v", got)
	}
}

func TestRunner_SucceedsOnFirstAttempt(t *testing.T) {
	p := backoff.DefaultPolicy()
	calls := 0
	err := p.Runner(3, noSleep, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRunner_RetriesOnFailure(t *testing.T) {
	p := backoff.DefaultPolicy()
	calls := 0
	sentinel := errors.New("transient")
	err := p.Runner(3, noSleep, func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRunner_ReturnsLastError(t *testing.T) {
	p := backoff.DefaultPolicy()
	sentinel := errors.New("permanent")
	err := p.Runner(3, noSleep, func() error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestRunner_SleepCalledBetweenAttempts(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0.0,
	}
	sleepCalls := 0
	_ = p.Runner(3, func(d time.Duration) { sleepCalls++ }, func() error {
		return errors.New("fail")
	})
	if sleepCalls != 2 {
		t.Fatalf("expected 2 sleep calls (between 3 attempts), got %d", sleepCalls)
	}
}
