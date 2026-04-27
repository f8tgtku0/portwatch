package jitter_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/jitter"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func openedMsg() notify.Message {
	return notify.Message{
		Change: state.Change{Port: 8080, Action: state.Opened},
	}
}

func TestNew_ZeroMax_CallsNextImmediately(t *testing.T) {
	j := jitter.New(0)
	called := false
	err := j.Delay(context.Background(), openedMsg(), func(_ notify.Message) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next to be called")
	}
}

func TestNew_PositiveMax_EventuallyCallsNext(t *testing.T) {
	j := jitter.New(10 * time.Millisecond)
	called := false
	err := j.Delay(context.Background(), openedMsg(), func(_ notify.Message) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next to be called after delay")
	}
}

func TestDelay_CancelledContext_DoesNotCallNext(t *testing.T) {
	j := jitter.New(5 * time.Second) // long enough that cancel fires first
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var called atomic.Bool
	_ = j.Delay(ctx, openedMsg(), func(_ notify.Message) error {
		called.Store(true)
		return nil
	})
	if called.Load() {
		t.Fatal("next should not be called when context is cancelled")
	}
}

func TestDelay_PropagatesNextError(t *testing.T) {
	j := jitter.New(0)
	sentinel := errors.New("downstream error")
	err := j.Delay(context.Background(), openedMsg(), func(_ notify.Message) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestMiddleware_Send_ForwardsMessage(t *testing.T) {
	var received notify.Message
	stub := &stubNotifier{fn: func(m notify.Message) error {
		received = m
		return nil
	}}
	mw := jitter.NewMiddleware(0, stub)
	msg := openedMsg()
	if err := mw.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Change.Port != msg.Change.Port {
		t.Fatalf("expected port %d, got %d", msg.Change.Port, received.Change.Port)
	}
}

type stubNotifier struct {
	fn func(notify.Message) error
}

func (s *stubNotifier) Send(m notify.Message) error { return s.fn(m) }
