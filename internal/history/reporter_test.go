package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleEntries() []Entry {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base, Port: 8080, Event: "opened"},
		{Timestamp: base.Add(time.Minute), Port: 3000, Event: "closed"},
	}
}

func TestPrint_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	if err := r.Print(sampleEntries()); err != nil {
		t.Fatalf("Print: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"TIMESTAMP", "PORT", "EVENT"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing header %q", want)
		}
	}
}

func TestPrint_ContainsEntryData(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	if err := r.Print(sampleEntries()); err != nil {
		t.Fatalf("Print: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "8080") || !strings.Contains(out, "opened") {
		t.Errorf("output missing entry data: %s", out)
	}
	if !strings.Contains(out, "3000") || !strings.Contains(out, "closed") {
		t.Errorf("output missing entry data: %s", out)
	}
}

func TestPrint_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	if err := r.Print(nil); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestSummary_Counts(t *testing.T) {
	s := Summary(sampleEntries())
	if !strings.Contains(s, "2 event(s)") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "1 opened") || !strings.Contains(s, "1 closed") {
		t.Errorf("unexpected counts in summary: %s", s)
	}
}

func TestSummary_Empty(t *testing.T) {
	s := Summary(nil)
	if !strings.Contains(s, "0 event(s)") {
		t.Errorf("unexpected summary for empty: %s", s)
	}
}
