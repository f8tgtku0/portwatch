package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
)

func fixedMsg() notify.Message {
	return notify.Message{
		Level:     notify.LevelAlert,
		Title:     "Port opened",
		Body:      "port 8080 is now open",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestLogNotifier_Send_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewLogNotifier(&buf)

	if err := n.Send(fixedMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ALERT", "Port opened", "port 8080 is now open", "2024-01-15"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q missing %q", out, want)
		}
	}
}

func TestLogNotifier_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil.
	n := notify.NewLogNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestLogNotifier_ZeroTimestamp_UsesNow(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewLogNotifier(&buf)
	msg := notify.Message{Level: notify.LevelInfo, Title: "t", Body: "b"}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

// errNotifier always returns an error.
type errNotifier struct{}

func (e *errNotifier) Send(_ notify.Message) error {
	return errors.New("backend unavailable")
}

func TestMulti_Send_DeliveresToAll(t *testing.T) {
	var b1, b2 bytes.Buffer
	m := notify.NewMulti(
		notify.NewLogNotifier(&b1),
		notify.NewLogNotifier(&b2),
	)
	if err := m.Send(fixedMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b1.Len() == 0 || b2.Len() == 0 {
		t.Error("expected both notifiers to receive the message")
	}
}

func TestMulti_Send_ReturnsFirstError(t *testing.T) {
	var buf bytes.Buffer
	m := notify.NewMulti(&errNotifier{}, notify.NewLogNotifier(&buf))
	err := m.Send(fixedMsg())
	if err == nil {
		t.Fatal("expected an error from errNotifier")
	}
	// Second notifier should still have been called.
	if buf.Len() == 0 {
		t.Error("expected second notifier to still receive the message")
	}
}
