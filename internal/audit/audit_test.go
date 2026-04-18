package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestRecord_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	if err := l.Record("port.opened", 8080, "new listener"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Event != "port.opened" {
		t.Errorf("expected event port.opened, got %q", e.Event)
	}
	if e.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.Port)
	}
	if e.Detail != "new listener" {
		t.Errorf("unexpected detail: %q", e.Detail)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.Record("port.opened", 22, "")
	_ = l.Record("port.closed", 22, "")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
