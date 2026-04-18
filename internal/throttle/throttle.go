// Package throttle provides a token-bucket style throttle for limiting
// how frequently change notifications are dispatched per port.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits calls per unique key to at most Burst events within Window.
type Throttle struct {
	mu     sync.Mutex
	window time.Duration
	burst  int
	buckets map[string][]time.Time
}

// New creates a Throttle allowing up to burst events per window per key.
func New(window time.Duration, burst int) *Throttle {
	if burst < 1 {
		burst = 1
	}
	return &Throttle{
		window:  window,
		burst:   burst,
		buckets: make(map[string][]time.Time),
	}
}

// Allow returns true if the event for key is within the allowed burst limit.
func (t *Throttle) Allow(key string) bool {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-t.window)
	times := t.buckets[key]

	// Prune old entries
	filtered := times[:0]
	for _, ts := range times {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= t.burst {
		t.buckets[key] = filtered
		return false
	}

	t.buckets[key] = append(filtered, now)
	return true
}

// Reset clears all tracking state.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buckets = make(map[string][]time.Time)
}

// Len returns the number of tracked keys.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.buckets)
}
