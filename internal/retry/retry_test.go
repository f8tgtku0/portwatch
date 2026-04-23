package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
)

type mockNotifier struct {
	calls   int
	failFor int // fail for the first N calls
	last    notify.Message
}

func (m *mockNotifier) Send(msg notify.Message) error {
	m.calls++
	m.last = msg
	if m.calls <= m.failFor {
		return errors.New("mock send error")
	}
	return nil
}

func noSleep(_ time.Duration) {}

func TestSend_SucceedsOnFirstAttempt(t *testing.T) {
	mock := &mockNotifier{}
	n := New(mock, DefaultPolicy())
	n.sleep = noSleep

	if err := n.Send(notify.Message{}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if mock.calls != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls)
	}
}

func TestSend_RetriesOnFailure(t *testing.T) {
	mock := &mockNotifier{failFor: 2}
	policy := Policy{MaxAttempts: 3, Delay: 0}
	n := New(mock, policy)
	n.sleep = noSleep

	if err := n.Send(notify.Message{}); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if mock.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", mock.calls)
	}
}

func TestSend_ReturnsErrorAfterAllAttempts(t *testing.T) {
	mock := &mockNotifier{failFor: 10}
	policy := Policy{MaxAttempts: 3, Delay: 0}
	n := New(mock, policy)
	n.sleep = noSleep

	err := n.Send(notify.Message{})
	if err == nil {
		t.Fatal("expected error after all attempts failed")
	}
	if mock.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", mock.calls)
	}
}

func TestSend_NilInner_ReturnsNil(t *testing.T) {
	n := New(nil, DefaultPolicy())
	n.sleep = noSleep
	if err := n.Send(notify.Message{}); err != nil {
		t.Fatalf("expected nil for nil inner, got %v", err)
	}
}

func TestSend_SleepCalledBetweenAttempts(t *testing.T) {
	mock := &mockNotifier{failFor: 2}
	policy := Policy{MaxAttempts: 3, Delay: 100 * time.Millisecond}
	n := New(mock, policy)

	sleptCount := 0
	n.sleep = func(d time.Duration) {
		sleptCount++
		if d != 100*time.Millisecond {
			t.Errorf("expected delay %v, got %v", policy.Delay, d)
		}
	}

	_ = n.Send(notify.Message{})
	if sleptCount != 2 {
		t.Fatalf("expected 2 sleeps between 3 attempts, got %d", sleptCount)
	}
}

func TestDefaultPolicy_Values(t *testing.T) {
	p := DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Delay != 500*time.Millisecond {
		t.Errorf("expected Delay=500ms, got %v", p.Delay)
	}
}
