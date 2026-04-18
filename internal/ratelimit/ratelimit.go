// Package ratelimit provides alert rate limiting to suppress duplicate notifications.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alerts for the same port within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New returns a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event key has not been seen within the cooldown window.
// If allowed, the key's timestamp is updated.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the rate limit record for a specific key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Clear removes all tracked keys.
func (l *Limiter) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}

// Len returns the number of tracked keys.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.last)
}
