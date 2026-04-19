// Package cooldown provides per-port alert cooldown tracking to prevent
// alert fatigue by suppressing repeated notifications within a configurable window.
package cooldown

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Cooldown tracks the last alert time per port+action key.
type Cooldown struct {
	mu      sync.Mutex
	last    map[string]time.Time
	window  time.Duration
}

// New creates a Cooldown with the given suppression window.
func New(window time.Duration) *Cooldown {
	return &Cooldown{
		last:   make(map[string]time.Time),
		window: window,
	}
}

// Allow returns true if enough time has passed since the last alert for the
// given port and action. It updates the last-seen timestamp on success.
func (c *Cooldown) Allow(port int, action state.ChangeType) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%d:%s", port, action)
	now := time.Now()

	if t, ok := c.last[key]; ok && now.Sub(t) < c.window {
		return false
	}

	c.last[key] = now
	return true
}

// Reset clears all cooldown state.
func (c *Cooldown) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last = make(map[string]time.Time)
}

// Prune removes entries older than the cooldown window.
func (c *Cooldown) Prune() {
	c.mu.Lock()
	defer c.mu.Unlock()

	cutoff := time.Now().Add(-c.window)
	for k, t := range c.last {
		if t.Before(cutoff) {
			delete(c.last, k)
		}
	}
}

// Len returns the number of tracked keys.
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.last)
}
