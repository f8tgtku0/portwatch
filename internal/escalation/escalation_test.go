package escalation_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/escalation"
	"github.com/user/portwatch/internal/state"
)

type fakeNotifier struct {
	count int
}

func (f *fakeNotifier) Send(_ state.PortChange) error {
	f.count++
	return nil
}

func TestTrack_AddsPending(t *testing.T) {
	esc := escalation.New(nil)
	esc.Track("opened:8080")
	if esc.PendingCount() != 1 {
		t.Fatalf("expected 1 pending, got %d", esc.PendingCount())
	}
}

func TestAcknowledge_RemovesPending(t *testing.T) {
	esc := escalation.New(nil)
	esc.Track("opened:8080")
	esc.Acknowledge("opened:8080")
	if esc.PendingCount() != 0 {
		t.Fatalf("expected 0 pending, got %d", esc.PendingCount())
	}
}

func TestEvaluate_FiresTierAfterDuration(t *testing.T) {
	n := &fakeNotifier{}
	esc := escalation.New([]escalation.Tier{
		{After: 0, Notifier: n},
	})
	esc.Track("opened:9090")
	esc.Evaluate(state.PortChange{Port: 9090, Action: state.Opened})
	if n.count != 1 {
		t.Fatalf("expected notifier called once, got %d", n.count)
	}
}

func TestEvaluate_DoesNotFireTierTwice(t *testing.T) {
	n := &fakeNotifier{}
	esc := escalation.New([]escalation.Tier{
		{After: 0, Notifier: n},
	})
	esc.Track("opened:9090")
	esc.Evaluate(state.PortChange{Port: 9090, Action: state.Opened})
	esc.Evaluate(state.PortChange{Port: 9090, Action: state.Opened})
	if n.count != 1 {
		t.Fatalf("expected notifier called once, got %d", n.count)
	}
}

func TestEvaluate_DoesNotFireBeforeDuration(t *testing.T) {
	n := &fakeNotifier{}
	esc := escalation.New([]escalation.Tier{
		{After: 10 * time.Hour, Notifier: n},
	})
	esc.Track("opened:7070")
	esc.Evaluate(state.PortChange{Port: 7070, Action: state.Opened})
	if n.count != 0 {
		t.Fatalf("expected notifier not called, got %d", n.count)
	}
}

func TestTrack_Idempotent(t *testing.T) {
	esc := escalation.New(nil)
	esc.Track("opened:8080")
	esc.Track("opened:8080")
	if esc.PendingCount() != 1 {
		t.Fatalf("expected 1 pending, got %d", esc.PendingCount())
	}
}
