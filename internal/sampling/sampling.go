// Package sampling provides probabilistic sampling for port change notifications.
// It allows reducing alert volume by forwarding only a fraction of changes,
// useful in high-churn environments where not every event needs to be actioned.
package sampling

import (
	"math/rand"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Sampler decides whether a given change should be forwarded based on a
// configured sampling rate in the range (0.0, 1.0].
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New creates a Sampler with the given rate. A rate of 1.0 forwards all
// changes; 0.5 forwards roughly half. Rates outside (0, 1] are clamped.
func New(rate float64) *Sampler {
	if rate <= 0 {
		rate = 0.01
	}
	if rate > 1.0 {
		rate = 1.0
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Allow returns true if the change should be forwarded according to the
// configured sampling rate.
func (s *Sampler) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rng.Float64() < s.rate
}

// Sample filters the provided changes, returning only those that pass the
// sampling check. Order is preserved for changes that are kept.
func (s *Sampler) Sample(changes []state.Change) []state.Change {
	if len(changes) == 0 {
		return nil
	}
	out := make([]state.Change, 0, len(changes))
	for _, c := range changes {
		if s.Allow() {
			out = append(out, c)
		}
	}
	return out
}

// Rate returns the current sampling rate.
func (s *Sampler) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}
