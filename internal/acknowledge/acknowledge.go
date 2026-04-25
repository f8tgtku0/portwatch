// Package acknowledge provides a mechanism for tracking and suppressing
// alerts that have been explicitly acknowledged by an operator.
package acknowledge

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Entry represents a single acknowledgement record.
type Entry struct {
	Port      int       `json:"port"`
	Action    string    `json:"action"`
	AckedAt   time.Time `json:"acked_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Note      string    `json:"note,omitempty"`
}

// Store holds acknowledged port change events.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised acknowledgement store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
	}
}

// Acknowledge marks a port/action pair as acknowledged for the given TTL.
func (s *Store) Acknowledge(port int, action, note string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	key := entryKey(port, action)
	s.entries[key] = Entry{
		Port:      port,
		Action:    action,
		AckedAt:   now,
		ExpiresAt: now.Add(ttl),
		Note:      note,
	}
}

// IsAcknowledged reports whether the given port/action pair is currently
// acknowledged (i.e. the acknowledgement exists and has not expired).
func (s *Store) IsAcknowledged(port int, action string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[entryKey(port, action)]
	if !ok {
		return false
	}
	return time.Now().Before(e.ExpiresAt)
}

// Revoke removes an acknowledgement before it naturally expires.
func (s *Store) Revoke(port int, action string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, entryKey(port, action))
}

// Prune removes all expired acknowledgements from the store.
func (s *Store) Prune() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, e := range s.entries {
		if now.After(e.ExpiresAt) {
			delete(s.entries, k)
		}
	}
}

// Active returns a copy of all non-expired acknowledgement entries.
func (s *Store) Active() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		if now.Before(e.ExpiresAt) {
			out = append(out, e)
		}
	}
	return out
}

func entryKey(port int, action string) string {
	_ = state.New // ensure import is used only when needed
	return action + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
