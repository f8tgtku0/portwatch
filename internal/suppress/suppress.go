// Package suppress provides a time-based suppression window for alerts,
// allowing operators to silence notifications during maintenance periods.
package suppress

import (
	"sync"
	"time"
)

// Window represents an active suppression period.
type Window struct {
	Start  time.Time
	End    time.Time
	Reason string
}

// Suppressor manages suppression windows.
type Suppressor struct {
	mu      sync.RWMutex
	windows []Window
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{}
}

// Add registers a suppression window.
func (s *Suppressor) Add(start, end time.Time, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows = append(s.windows, Window{Start: start, End: end, Reason: reason})
}

// IsSuppressed reports whether the given time falls within any active window.
func (s *Suppressor) IsSuppressed(t time.Time) (bool, Window) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.windows {
		if !t.Before(w.Start) && t.Before(w.End) {
			return true, w
		}
	}
	return false, Window{}
}

// Prune removes windows that have already ended relative to now.
func (s *Suppressor) Prune(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	active := s.windows[:0]
	for _, w := range s.windows {
		if now.Before(w.End) {
			active = append(active, w)
		}
	}
	s.windows = active
}

// Active returns a copy of all current windows.
func (s *Suppressor) Active() []Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Window, len(s.windows))
	copy(out, s.windows)
	return out
}
