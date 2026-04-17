package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(Snapshot{})
	out := buf.String()
	for _, h := range []string{"METRIC", "VALUE"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestPrint_ContainsMetricValues(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	s := Snapshot{
		ScansTotal:   42,
		ChangesTotal: 7,
		AlertsTotal:  3,
		OpenPorts:    12,
		Uptime:       5 * time.Minute,
		LastScanAt:   time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}
	r.Print(s)
	out := buf.String()
	for _, want := range []string{"42", "7", "3", "12", "2024-01-15"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestPrint_NoLastScan_ShowsNever(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(Snapshot{})
	if !strings.Contains(buf.String(), "never") {
		t.Error("expected 'never' when LastScanAt is zero")
	}
}

func TestNewReporter_NilWriter_DefaultsToStdout(t *testing.T) {
	r := NewReporter(nil)
	if r.w == nil {
		t.Error("expected non-nil writer")
	}
}
