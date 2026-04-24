// Package window provides a time-windowed event counter that tracks how many
// port change events have occurred within a rolling duration. It can be used
// to detect burst activity and suppress or escalate alerts accordingly.
package window

import (
	"sync"
	"time"
)

// Counter tracks events within a sliding time window.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	events   []time.Time
	maxSize  int
}

// New creates a Counter with the given sliding window duration and optional
// maximum event history size. If maxSize is 0 it defaults to 1000.
func New(window time.Duration, maxSize int) *Counter {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &Counter{
		window:  window,
		maxSize: maxSize,
	}
}

// Record adds the current timestamp to the event log.
func (c *Counter) Record() {
	c.RecordAt(time.Now())
}

// RecordAt adds a specific timestamp — useful for testing.
func (c *Counter) RecordAt(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prune(t)
	if len(c.events) < c.maxSize {
		c.events = append(c.events, t)
	}
}

// Count returns the number of events within the current window.
func (c *Counter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prune(time.Now())
	return len(c.events)
}

// Exceeds reports whether the event count within the window exceeds threshold.
func (c *Counter) Exceeds(threshold int) bool {
	return c.Count() > threshold
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}

// prune removes events older than the window. Must be called with mu held.
func (c *Counter) prune(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.events) && c.events[i].Before(cutoff) {
		i++
	}
	c.events = c.events[i:]
}
