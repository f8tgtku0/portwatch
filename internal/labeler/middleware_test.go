package labeler_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/labeler"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

type captureNotifier struct {
	last notify.Message
	err  error
}

func (c *captureNotifier) Send(msg notify.Message) error {
	c.last = msg
	return c.err
}

func makeChanges(ports ...int) []state.Change {
	out := make([]state.Change, len(ports))
	for i, p := range ports {
		out[i] = state.Change{Port: p, Action: state.Opened}
	}
	return out
}

func TestMiddleware_Send_LabelsChanges(t *testing.T) {
	cap := &captureNotifier{}
	l := labeler.New(nil)
	mw := labeler.NewMiddleware(l, cap)

	msg := notify.Message{Changes: makeChanges(22, 80)}
	if err := mw.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.last.Changes[0].Label != "ssh" {
		t.Errorf("port 22: expected ssh, got %q", cap.last.Changes[0].Label)
	}
	if cap.last.Changes[1].Label != "http" {
		t.Errorf("port 80: expected http, got %q", cap.last.Changes[1].Label)
	}
}

func TestMiddleware_Send_NilLabeler_ForwardsUnchanged(t *testing.T) {
	cap := &captureNotifier{}
	mw := labeler.NewMiddleware(nil, cap)

	changes := makeChanges(22)
	msg := notify.Message{Changes: changes}
	if err := mw.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.last.Changes[0].Label != "" {
		t.Errorf("expected empty label, got %q", cap.last.Changes[0].Label)
	}
}

func TestMiddleware_Send_PropagatesError(t *testing.T) {
	sentinel := errors.New("downstream error")
	cap := &captureNotifier{err: sentinel}
	l := labeler.New(nil)
	mw := labeler.NewMiddleware(l, cap)

	if err := mw.Send(notify.Message{Changes: makeChanges(80)}); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestMiddleware_Apply_AnnotatesChanges(t *testing.T) {
	l := labeler.New(map[int]string{5432: "postgres"})
	mw := labeler.NewMiddleware(l, &captureNotifier{})

	out := mw.Apply(makeChanges(5432))
	if out[0].Label != "postgres" {
		t.Errorf("expected postgres, got %q", out[0].Label)
	}
}
