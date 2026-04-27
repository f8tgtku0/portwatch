package replay_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/replay"
)

type stubNotifier struct {
	received []notify.Message
	errOn    int
	calls    int
}

func (s *stubNotifier) Send(msg notify.Message) error {
	s.calls++
	if s.errOn > 0 && s.calls == s.errOn {
		return errors.New("stub error")
	}
	s.received = append(s.received, msg)
	return nil
}

func openedMsg(port int) notify.Message {
	return notify.Message{Port: port, Action: "opened", At: time.Now()}
}

func TestRecord_BuffersEntry(t *testing.T) {
	r := replay.New(10, nil)
	r.Record(openedMsg(8080))
	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	r := replay.New(3, nil)
	for i := 0; i < 5; i++ {
		r.Record(openedMsg(8080 + i))
	}
	if r.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", r.Len())
	}
}

func TestReplay_DeliversAllEntries(t *testing.T) {
	r := replay.New(10, nil)
	r.Record(openedMsg(80))
	r.Record(openedMsg(443))

	stub := &stubNotifier{}
	delivered := r.Replay(stub)

	if delivered != 2 {
		t.Fatalf("expected 2 delivered, got %d", delivered)
	}
	if len(stub.received) != 2 {
		t.Fatalf("expected 2 received messages")
	}
}

func TestReplay_LogsErrorAndContinues(t *testing.T) {
	var buf bytes.Buffer
	r := replay.New(10, &buf)
	r.Record(openedMsg(80))
	r.Record(openedMsg(443))

	stub := &stubNotifier{errOn: 1}
	delivered := r.Replay(stub)

	if delivered != 1 {
		t.Fatalf("expected 1 delivered after one error, got %d", delivered)
	}
	if buf.Len() == 0 {
		t.Fatal("expected error to be logged")
	}
}

func TestClear_EmptiesBuffer(t *testing.T) {
	r := replay.New(10, nil)
	r.Record(openedMsg(22))
	r.Clear()
	if r.Len() != 0 {
		t.Fatalf("expected 0 entries after clear, got %d", r.Len())
	}
}
