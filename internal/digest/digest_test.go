package digest_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/state"
)

func openedPorts(nums ...int) []state.Port {
	var out []state.Port
	for _, n := range nums {
		out = append(out, state.Port{Number: n})
	}
	return out
}

func TestRecord_AccumulatesEntries(t *testing.T) {
	d := digest.New(nil)
	d.Record(openedPorts(80, 443), openedPorts(8080))
	if d.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", d.Len())
	}
}

func TestFlush_WritesAndClears(t *testing.T) {
	var buf bytes.Buffer
	d := digest.New(&buf)
	d.Record(openedPorts(22), openedPorts(9000))
	d.Flush()
	out := buf.String()
	if !strings.Contains(out, "1 opened") {
		t.Errorf("expected opened count in output, got: %s", out)
	}
	if !strings.Contains(out, "1 closed") {
		t.Errorf("expected closed count in output, got: %s", out)
	}
	if d.Len() != 0 {
		t.Errorf("expected entries cleared after flush, got %d", d.Len())
	}
}

func TestFlush_EmptyIsNoop(t *testing.T) {
	var buf bytes.Buffer
	d := digest.New(&buf)
	d.Flush()
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty digest, got: %s", buf.String())
	}
}

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	d := digest.New(nil)
	if d == nil {
		t.Fatal("expected non-nil digest")
	}
}

func TestFlush_ContainsPortNumbers(t *testing.T) {
	var buf bytes.Buffer
	d := digest.New(&buf)
	d.Record(openedPorts(3306), nil)
	d.Flush()
	if !strings.Contains(buf.String(), "3306") {
		t.Errorf("expected port 3306 in output, got: %s", buf.String())
	}
}
