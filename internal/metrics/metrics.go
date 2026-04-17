// Package metrics tracks runtime statistics for portwatch.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	ScansTotal   int64
	ChangesTotal int64
	AlertsTotal  int64
	LastScanAt   time.Time
	Uptime       time.Duration
	OpenPorts    int
}

// Collector accumulates runtime metrics in a thread-safe manner.
type Collector struct {
	mu           sync.RWMutex
	scansTotal   int64
	changesTotal int64
	alertsTotal  int64
	lastScanAt   time.Time
	startedAt    time.Time
	openPorts    int
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{startedAt: time.Now()}
}

// RecordScan increments the scan counter and records open port count.
func (c *Collector) RecordScan(openPorts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
	c.openPorts = openPorts
}

// RecordChange increments the change counter by delta.
func (c *Collector) RecordChange(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changesTotal += int64(delta)
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// Snapshot returns a copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:   c.scansTotal,
		ChangesTotal: c.changesTotal,
		AlertsTotal:  c.alertsTotal,
		LastScanAt:   c.lastScanAt,
		Uptime:       time.Since(c.startedAt).Truncate(time.Second),
		OpenPorts:    c.openPorts,
	}
}
