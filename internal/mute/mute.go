// Package mute provides a time-bounded suppression mechanism that silences
// notifications for specific ports during scheduled maintenance windows or
// operator-defined quiet periods.
package mute

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Entry represents a single mute rule keyed by port number.
type Entry struct {
	Port      int
	Reason    string
	ExpiresAt time.Time
}

// Muter holds the set of currently active mute rules.
type Muter struct {
	mu      sync.RWMutex
	entries map[int]Entry
}

// New returns an initialised Muter with no active rules.
func New() *Muter {
	return &Muter{
		entries: make(map[int]Entry),
	}
}

// Add registers a mute rule for the given port that expires after ttl.
// Calling Add for an already-muted port replaces the existing rule.
func (m *Muter) Add(port int, reason string, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[port] = Entry{
		Port:      port,
		Reason:    reason,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Remove deletes the mute rule for the given port immediately.
func (m *Muter) Remove(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, port)
}

// IsMuted reports whether the given port is currently muted.
// Expired rules are treated as inactive and pruned lazily.
func (m *Muter) IsMuted(port int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[port]
	if !ok {
		return false
	}
	if time.Now().After(e.ExpiresAt) {
		delete(m.entries, port)
		return false
	}
	return true
}

// Active returns a snapshot of all currently active (non-expired) mute rules.
func (m *Muter) Active() []Entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	out := make([]Entry, 0, len(m.entries))
	for port, e := range m.entries {
		if now.After(e.ExpiresAt) {
			delete(m.entries, port)
			continue
		}
		out = append(out, e)
	}
	return out
}

// Prune removes all expired entries from the internal map.
func (m *Muter) Prune() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for port, e := range m.entries {
		if now.After(e.ExpiresAt) {
			delete(m.entries, port)
		}
	}
}

// Ensure state import is used only when needed; kept for future persistence.
var _ = state.New
