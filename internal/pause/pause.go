// Package pause provides a mechanism to temporarily suspend alert delivery
// without stopping the underlying port scanner. While paused, changes are
// still detected but notifications are silently dropped.
package pause

import (
	"sync"
	"time"
)

// Pauser tracks whether alerting is currently paused and when the pause expires.
type Pauser struct {
	mu       sync.RWMutex
	paused   bool
	until    time.Time
	resumeAt func() time.Time // injectable for testing
}

// New returns a new Pauser that is initially active (not paused).
func New() *Pauser {
	return &Pauser{
		resumeAt: time.Now,
	}
}

// Pause suspends alert delivery for the given duration.
// Calling Pause while already paused extends the deadline.
func (p *Pauser) Pause(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = true
	p.until = p.resumeAt().Add(d)
}

// Resume cancels an active pause immediately.
func (p *Pauser) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
	p.until = time.Time{}
}

// IsPaused reports whether alerting is currently suspended.
// An expired timed pause is automatically cleared.
func (p *Pauser) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.paused {
		return false
	}
	if !p.until.IsZero() && p.resumeAt().After(p.until) {
		p.paused = false
		p.until = time.Time{}
		return false
	}
	return true
}

// Until returns the time at which the current pause will expire.
// Returns the zero value if not paused or paused indefinitely.
func (p *Pauser) Until() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.until
}
