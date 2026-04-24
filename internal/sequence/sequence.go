// Package sequence assigns monotonically increasing sequence numbers to
// change events, enabling downstream consumers to detect gaps or reordering.
package sequence

import (
	"sync"

	"portwatch/internal/state"
)

// Sequencer stamps each batch of changes with a sequence number.
type Sequencer struct {
	mu      sync.Mutex
	counter uint64
}

// Stamped wraps a change with its assigned sequence number.
type Stamped struct {
	Seq    uint64
	Change state.Change
}

// New returns a new Sequencer starting at sequence number 1.
func New() *Sequencer {
	return &Sequencer{}
}

// Stamp assigns the next sequence number to each change in the slice and
// returns the annotated results. Each call increments the counter by one,
// shared across all changes in that batch.
func (s *Sequencer) Stamp(changes []state.Change) []Stamped {
	if len(changes) == 0 {
		return nil
	}

	s.mu.Lock()
	s.counter++
	seq := s.counter
	s.mu.Unlock()

	out := make([]Stamped, len(changes))
	for i, c := range changes {
		out[i] = Stamped{Seq: seq, Change: c}
	}
	return out
}

// Current returns the most recently issued sequence number without
// incrementing the counter.
func (s *Sequencer) Current() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counter
}

// Reset sets the counter back to zero. Intended for testing only.
func (s *Sequencer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = 0
}
