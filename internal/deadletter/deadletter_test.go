package deadletter_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/deadletter"
	"github.com/user/portwatch/internal/notify"
)

type failNotifier struct{ err error }

func (f *failNotifier) Send(_ notify.Message) error { return f.err }

type okNotifier struct{ called bool }

func (o *okNotifier) Send(_ notify.Message) error { o.called = true; return nil }

func sampleMsg(port int) notify.Message {
	return notify.Message{Port: port, Action: "opened"}
}

func TestRecord_AddsEntry(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	q.Record(sampleMsg(8080), "timeout")
	if q.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", q.Len())
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	q := deadletter.New(3, &bytes.Buffer{})
	for i := 0; i < 4; i++ {
		q.Record(sampleMsg(8080+i), "err")
	}
	if q.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", q.Len())
	}
	entries := q.Drain()
	if entries[0].Message.Port != 8081 {
		t.Errorf("expected oldest evicted; first port = %d", entries[0].Message.Port)
	}
}

func TestDrain_ClearsQueue(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	q.Record(sampleMsg(443), "refused")
	q.Drain()
	if q.Len() != 0 {
		t.Fatalf("expected empty queue after drain")
	}
}

func TestRecord_TimestampIsRecent(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	before := time.Now()
	q.Record(sampleMsg(22), "err")
	after := time.Now()
	entries := q.Drain()
	if entries[0].At.Before(before) || entries[0].At.After(after) {
		t.Errorf("timestamp out of range: %v", entries[0].At)
	}
}

func TestMiddleware_OnSuccess_DoesNotQueue(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	ok := &okNotifier{}
	mw := deadletter.NewMiddleware(ok, q)
	if err := mw.Send(sampleMsg(80)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Len() != 0 {
		t.Errorf("expected empty queue on success")
	}
}

func TestMiddleware_OnFailure_RecordsAndReturnsError(t *testing.T) {
	q := deadletter.New(10, &bytes.Buffer{})
	sentinel := errors.New("dial timeout")
	mw := deadletter.NewMiddleware(&failNotifier{err: sentinel}, q)
	err := mw.Send(sampleMsg(9000))
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if q.Len() != 1 {
		t.Errorf("expected 1 dead-letter entry, got %d", q.Len())
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	q := deadletter.New(0, &bytes.Buffer{})
	for i := 0; i < 150; i++ {
		q.Record(sampleMsg(i), "err")
	}
	if q.Len() != 100 {
		t.Errorf("expected default cap 100, got %d", q.Len())
	}
}
