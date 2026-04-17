package metrics

import (
	"testing"
	"time"
)

func TestNew_InitializesStartTime(t *testing.T) {
	before := time.Now()
	c := New()
	after := time.Now()
	if c.startedAt.Before(before) || c.startedAt.After(after) {
		t.Error("startedAt not within expected range")
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := New()
	c.RecordScan(5)
	c.RecordScan(3)
	s := c.Snapshot()
	if s.ScansTotal != 2 {
		t.Errorf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.OpenPorts != 3 {
		t.Errorf("expected 3 open ports, got %d", s.OpenPorts)
	}
}

func TestRecordScan_UpdatesLastScanAt(t *testing.T) {
	c := New()
	before := time.Now()
	c.RecordScan(0)
	after := time.Now()
	s := c.Snapshot()
	if s.LastScanAt.Before(before) || s.LastScanAt.After(after) {
		t.Error("LastScanAt not within expected range")
	}
}

func TestRecordChange_AccumulatesDelta(t *testing.T) {
	c := New()
	c.RecordChange(2)
	c.RecordChange(3)
	if c.Snapshot().ChangesTotal != 5 {
		t.Errorf("expected 5 changes, got %d", c.Snapshot().ChangesTotal)
	}
}

func TestRecordAlert_IncrementsCounter(t *testing.T) {
	c := New()
	c.RecordAlert()
	c.RecordAlert()
	if c.Snapshot().AlertsTotal != 2 {
		t.Errorf("expected 2 alerts, got %d", c.Snapshot().AlertsTotal)
	}
}

func TestSnapshot_UptimeGrowsOverTime(t *testing.T) {
	c := New()
	time.Sleep(10 * time.Millisecond)
	if c.Snapshot().Uptime <= 0 {
		t.Error("expected positive uptime")
	}
}
