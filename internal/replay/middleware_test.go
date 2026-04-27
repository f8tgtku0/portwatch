package replay_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/replay"
)

type failNotifier struct{}

func (f *failNotifier) Send(_ notify.Message) error {
	return errors.New("downstream failure")
}

func TestMiddleware_RecordsOnSuccess(t *testing.T) {
	r := replay.New(10, nil)
	stub := &stubNotifier{}
	mw := replay.NewMiddleware(stub, r)

	if err := mw.Send(openedMsg(8080)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 buffered entry, got %d", r.Len())
	}
}

func TestMiddleware_DoesNotRecordOnError(t *testing.T) {
	r := replay.New(10, nil)
	mw := replay.NewMiddleware(&failNotifier{}, r)

	_ = mw.Send(openedMsg(9090))

	if r.Len() != 0 {
		t.Fatalf("expected 0 buffered entries on downstream error, got %d", r.Len())
	}
}

func TestMiddleware_NilReplayer_PassesThrough(t *testing.T) {
	stub := &stubNotifier{}
	mw := replay.NewMiddleware(stub, nil)

	if err := mw.Send(openedMsg(22)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stub.received) != 1 {
		t.Fatalf("expected message forwarded to underlying notifier")
	}
}

func TestMiddleware_PropagatesDownstreamError(t *testing.T) {
	r := replay.New(10, nil)
	mw := replay.NewMiddleware(&failNotifier{}, r)

	err := mw.Send(openedMsg(443))
	if err == nil {
		t.Fatal("expected error to propagate from downstream notifier")
	}
}
