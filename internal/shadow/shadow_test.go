package shadow_test

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/shadow"
)

type recordNotifier struct {
	mu   sync.Mutex
	msgs []notify.Message
	err  error
}

func (r *recordNotifier) Send(msg notify.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.msgs = append(r.msgs, msg)
	return r.err
}

func TestShadow_PrimaryReceivesMessage(t *testing.T) {
	primary := &recordNotifier{}
	shadowN := &recordNotifier{}
	n := shadow.New(primary, shadowN, nil)

	msg := notify.Message{Text: "port 8080 opened"}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(primary.msgs) != 1 {
		t.Fatalf("expected 1 primary message, got %d", len(primary.msgs))
	}
}

func TestShadow_ShadowReceivesMessage(t *testing.T) {
	primary := &recordNotifier{}
	shadowN := &recordNotifier{}
	n := shadow.New(primary, shadowN, nil)

	_ = n.Send(notify.Message{Text: "port 443 opened"})
	time.Sleep(20 * time.Millisecond)

	shadowN.mu.Lock()
	defer shadowN.mu.Unlock()
	if len(shadowN.msgs) != 1 {
		t.Fatalf("expected 1 shadow message, got %d", len(shadowN.msgs))
	}
}

func TestShadow_ShadowErrorDoesNotPropagate(t *testing.T) {
	primary := &recordNotifier{}
	shadowN := &recordNotifier{err: errors.New("shadow down")}
	var buf bytes.Buffer
	n := shadow.New(primary, shadowN, &buf)

	if err := n.Send(notify.Message{Text: "test"}); err != nil {
		t.Fatalf("primary error should be nil, got %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if buf.Len() == 0 {
		t.Error("expected shadow error to be logged")
	}
}

func TestShadow_PrimaryErrorPropagates(t *testing.T) {
	primary := &recordNotifier{err: errors.New("primary down")}
	shadowN := &recordNotifier{}
	n := shadow.New(primary, shadowN, nil)

	if err := n.Send(notify.Message{Text: "test"}); err == nil {
		t.Fatal("expected primary error to propagate")
	}
}

func TestNew_NilLog_DefaultsToStderr(t *testing.T) {
	primary := &recordNotifier{}
	shadowN := &recordNotifier{}
	n := shadow.New(primary, shadowN, nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
