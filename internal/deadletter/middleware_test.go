package deadletter_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/deadletter"
	"github.com/user/portwatch/internal/notify"
)

func TestMiddleware_ReasonContainsErrorText(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	mw := deadletter.NewMiddleware(&failNotifier{err: errors.New("connection refused")}, q)
	_ = mw.Send(notify.Message{Port: 3306, Action: "opened"})
	entries := q.Drain()
	if len(entries) == 0 {
		t.Fatal("expected entry in dead-letter queue")
	}
	if entries[0].Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestMiddleware_ImplementsNotifier(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	var _ notify.Notifier = deadletter.NewMiddleware(&okNotifier{}, q)
}

func TestMiddleware_MultipleFailures_AllQueued(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	mw := deadletter.NewMiddleware(&failNotifier{err: errors.New("err")}, q)
	for _, port := range []int{80, 443, 8080} {
		_ = mw.Send(notify.Message{Port: port, Action: "closed"})
	}
	if q.Len() != 3 {
		t.Errorf("expected 3 dead-letter entries, got %d", q.Len())
	}
}

func TestMiddleware_MessagePreservedInQueue(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	msg := notify.Message{Port: 5432, Action: "opened"}
	mw := deadletter.NewMiddleware(&failNotifier{err: errors.New("refused")}, q)
	_ = mw.Send(msg)
	entries := q.Drain()
	if entries[0].Message.Port != 5432 {
		t.Errorf("expected port 5432, got %d", entries[0].Message.Port)
	}
	if entries[0].Message.Action != "opened" {
		t.Errorf("expected action 'opened', got %s", entries[0].Message.Action)
	}
}
