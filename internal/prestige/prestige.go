// Package prestige tracks how many times a port has been seen opened
// across scan cycles, producing a "prestige" score that reflects
// long-term stability vs. ephemeral appearances.
package prestige

import (
	"sync"

	"github.com/user/portwatch/internal/state"
)

// Score holds the open-count and closed-count for a single port.
type Score struct {
	Port    int
	Opens   int
	Closes  int
}

// Tracker accumulates open/close events per port.
type Tracker struct {
	mu     sync.RWMutex
	scores map[int]*Score
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{scores: make(map[int]*Score)}
}

// Record ingests a slice of state.Change values and increments the
// appropriate counters for each port.
func (t *Tracker) Record(changes []state.Change) {
	if len(changes) == 0 {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, c := range changes {
		s, ok := t.scores[c.Port]
		if !ok {
			s = &Score{Port: c.Port}
			t.scores[c.Port] = s
		}
		if c.Action == state.Opened {
			s.Opens++
		} else {
			s.Closes++
		}
	}
}

// Get returns the Score for a port. The second return value is false
// when the port has never been recorded.
func (t *Tracker) Get(port int) (Score, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s, ok := t.scores[port]
	if !ok {
		return Score{}, false
	}
	return *s, true
}

// All returns a copy of every tracked Score.
func (t *Tracker) All() []Score {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Score, 0, len(t.scores))
	for _, s := range t.scores {
		out = append(out, *s)
	}
	return out
}

// Reset clears all accumulated scores.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.scores = make(map[int]*Score)
}
