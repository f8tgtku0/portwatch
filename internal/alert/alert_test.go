package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func TestNotify_OpenedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	err := n.Notify(monitor.Change{Port: 8080, Type: monitor.ChangeTypeOpened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
}

func TestNotify_ClosedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	err := n.Notify(monitor.Change{Port: 9090, Type: monitor.ChangeTypeClosed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level, got: %s", out)
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected port 9090 in output, got: %s", out)
	}
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
}

func TestNotifyAll_MultipleChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	changes := []monitor.Change{
		{Port: 80, Type: monitor.ChangeTypeOpened},
		{Port: 443, Type: monitor.ChangeTypeOpened},
		{Port: 8080, Type: monitor.ChangeTypeClosed},
	}

	if err := n.NotifyAll(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines of output, got %d:\n%s", len(lines), out)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
