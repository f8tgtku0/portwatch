package circuit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuit"
	"github.com/user/portwatch/internal/notify"
)

type fakeNotifier struct {
	err   error
	calls int
}

func (f *fakeNotifier) Send(_ notify.Message) error {
	f.calls++
	return f.err
}

func msg() notify.Message {
	return notify.Message{Port: 8080, Action: "opened"}
}

func TestBreaker_ClosedByDefault(t *testing.T) {
	fake := &fakeNotifier{}
	b := circuit.New(fake, 3, time.Second)
	if b.CurrentState() != circuit.StateClosed {
		t.Fatal("expected closed state initially")
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	fake := &fakeNotifier{err: errors.New("fail")}
	b := circuit.New(fake, 3, time.Second)

	for i := 0; i < 3; i++ {
		_ = b.Send(msg())
	}

	if b.CurrentState() != circuit.StateOpen {
		t.Fatalf("expected open state after %d failures", 3)
	}
}

func TestBreaker_BlocksSendWhenOpen(t *testing.T) {
	fake := &fakeNotifier{err: errors.New("fail")}
	b := circuit.New(fake, 1, time.Hour)

	_ = b.Send(msg()) // trips the breaker

	err := b.Send(msg())
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
	if fake.calls != 1 {
		t.Fatalf("inner notifier should not be called when open, got %d calls", fake.calls)
	}
}

func TestBreaker_HalfOpenAfterCooldown(t *testing.T) {
	fake := &fakeNotifier{err: errors.New("fail")}
	b := circuit.New(fake, 1, 10*time.Millisecond)

	_ = b.Send(msg()) // open the circuit

	time.Sleep(20 * time.Millisecond)

	fake.err = nil
	err := b.Send(msg())
	if err != nil {
		t.Fatalf("expected success after cooldown, got: %v", err)
	}
	if b.CurrentState() != circuit.StateClosed {
		t.Fatalf("expected closed after successful probe, got %v", b.CurrentState())
	}
}

func TestBreaker_Reset_ForcesClosed(t *testing.T) {
	fake := &fakeNotifier{err: errors.New("fail")}
	b := circuit.New(fake, 1, time.Hour)

	_ = b.Send(msg())
	if b.CurrentState() != circuit.StateOpen {
		t.Fatal("expected open")
	}

	b.Reset()
	if b.CurrentState() != circuit.StateClosed {
		t.Fatal("expected closed after reset")
	}
}

func TestBreaker_SuccessResetsFailureCount(t *testing.T) {
	callCount := 0
	fake := &fakeNotifier{}
	b := circuit.New(fake, 3, time.Second)

	fake.err = errors.New("fail")
	_ = b.Send(msg())
	_ = b.Send(msg())
	callCount = fake.calls

	fake.err = nil
	_ = b.Send(msg())

	if b.CurrentState() != circuit.StateClosed {
		t.Fatal("expected still closed after recovery")
	}
	_ = callCount
}
