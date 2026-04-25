package pause_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/pause"
	"github.com/user/portwatch/internal/state"
)

type recordingNotifier struct {
	called int
	err    error
}

func (r *recordingNotifier) Send(_ context.Context, _ notify.Message) error {
	r.called++
	return r.err
}

func TestMiddleware_ForwardsWhenNotPaused(t *testing.T) {
	rec := &recordingNotifier{}
	p := pause.New()
	mw := pause.NewMiddleware(p, rec)

	_ = mw.Send(context.Background(), notify.Message{})
	if rec.called != 1 {
		t.Fatalf("expected 1 call, got %d", rec.called)
	}
}

func TestMiddleware_DropsWhenPaused(t *testing.T) {
	rec := &recordingNotifier{}
	p := pause.New()
	p.Pause(5 * time.Minute)
	mw := pause.NewMiddleware(p, rec)

	err := mw.Send(context.Background(), notify.Message{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.called != 0 {
		t.Fatalf("expected 0 calls while paused, got %d", rec.called)
	}
}

func TestMiddleware_PropagatesError(t *testing.T) {
	sentinel := errors.New("send failed")
	rec := &recordingNotifier{err: sentinel}
	p := pause.New()
	mw := pause.NewMiddleware(p, rec)

	err := mw.Send(context.Background(), notify.Message{})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestMiddleware_Apply_ReturnsNilWhenPaused(t *testing.T) {
	rec := &recordingNotifier{}
	p := pause.New()
	p.Pause(5 * time.Minute)
	mw := pause.NewMiddleware(p, rec)

	changes := []state.Change{{Port: 80, Action: state.Opened}}
	got := mw.Apply(changes)
	if got != nil {
		t.Fatalf("expected nil when paused, got %v", got)
	}
}

func TestMiddleware_Apply_PassesThroughWhenNotPaused(t *testing.T) {
	rec := &recordingNotifier{}
	p := pause.New()
	mw := pause.NewMiddleware(p, rec)

	changes := []state.Change{{Port: 443, Action: state.Opened}}
	got := mw.Apply(changes)
	if len(got) != 1 {
		t.Fatalf("expected 1 change, got %d", len(got))
	}
}

func TestNewMiddleware_NilPauser_UsesDefault(t *testing.T) {
	rec := &recordingNotifier{}
	mw := pause.NewMiddleware(nil, rec)
	_ = mw.Send(context.Background(), notify.Message{})
	if rec.called != 1 {
		t.Fatalf("expected 1 call with nil pauser, got %d", rec.called)
	}
}
