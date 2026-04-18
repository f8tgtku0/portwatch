package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/state"
)

func TestAudit_RecordsOpenedAndClosed(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	ca := audit.NewChangeAuditor(l)

	changes := []state.Change{
		{Port: 80, Proto: "tcp", Opened: true},
		{Port: 443, Proto: "tcp", Opened: false},
	}
	out := ca.Audit(changes)
	if len(out) != 2 {
		t.Fatalf("expected 2 changes returned, got %d", len(out))
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 audit lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "port.opened") {
		t.Errorf("expected port.opened in first line, got %q", lines[0])
	}
	if !strings.Contains(lines[1], "port.closed") {
		t.Errorf("expected port.closed in second line, got %q", lines[1])
	}
}

func TestAudit_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	ca := audit.NewChangeAuditor(l)
	out := ca.Audit(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d", len(out))
	}
	if buf.Len() != 0 {
		t.Error("expected no output for empty changes")
	}
}
