// Package quota enforces per-channel notification quotas over a rolling window.
// When a channel exceeds its configured limit, further notifications are dropped
// until the window resets.
package quota

import (
	"sync"
	"time"
)

// Entry tracks the count and window start for a single channel key.
type Entry struct {
	count     int
	windowEnd time.Time
}

// Quota enforces a maximum number of notifications per channel per window.
type Quota struct {
	mu     sync.Mutex
	entries map[string]*Entry
	max    int
	window time.Duration
	now    func() time.Time
}

// New creates a Quota that allows at most max notifications per window duration.
func New(max int, window time.Duration) *Quota {
	return &Quota{
		entries: make(map[string]*Entry),
		max:    max,
		window: window,
		now:    time.Now,
	}
}

// Allow returns true if the channel identified by key is within quota.
// It increments the counter on each allowed call.
func (q *Quota) Allow(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	e, ok := q.entries[key]
	if !ok || now.After(e.windowEnd) {
		q.entries[key] = &Entry{count: 1, windowEnd: now.Add(q.window)}
		return true
	}
	if e.count >= q.max {
		return false
	}
	e.count++
	return true
}

// Remaining returns the number of notifications still allowed for key
// within the current window. Returns max if no window has started.
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	e, ok := q.entries[key]
	if !ok || q.now().After(e.windowEnd) {
		return q.max
	}
	rem := q.max - e.count
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the quota state for all keys.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries = make(map[string]*Entry)
}
