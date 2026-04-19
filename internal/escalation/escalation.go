// Package escalation provides multi-tier alert escalation based on how long
// a port change remains unacknowledged.
package escalation

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Tier defines a single escalation level.
type Tier struct {
	// After is the duration after the initial change before this tier fires.
	After time.Duration
	// Notifier is the notifier to invoke at this tier.
	Notifier notify.Notifier
}

// entry tracks an unacknowledged change.
type entry struct {
	firedTiers map[int]bool
	firstSeen  time.Time
}

// Escalator watches for unacknowledged changes and re-alerts on a schedule.
type Escalator struct {
	mu      sync.Mutex
	tiers   []Tier
	pending map[string]*entry // key -> entry
}

// New creates an Escalator with the given tiers.
func New(tiers []Tier) *Escalator {
	return &Escalator{
		tiers:   tiers,
		pending: make(map[string]*entry),
	}
}

// Track registers a change as unacknowledged.
func (e *Escalator) Track(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.pending[key]; !ok {
		e.pending[key] = &entry{
			firedTiers: make(map[int]bool),
			firstSeen:  time.Now(),
		}
	}
}

// Acknowledge removes a change from escalation tracking.
func (e *Escalator) Acknowledge(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.pending, key)
}

// Evaluate checks all pending changes and fires any due tiers.
func (e *Escalator) Evaluate(change state.PortChange) {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	for key, ent := range e.pending {
		for i, tier := range e.tiers {
			if ent.firedTiers[i] {
				continue
			}
			if now.Sub(ent.firstSeen) >= tier.After {
				_ = tier.Notifier.Send(change)
				ent.firedTiers[i] = true
				e.pending[key] = ent
			}
		}
	}
}

// PendingCount returns the number of unacknowledged changes.
func (e *Escalator) PendingCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.pending)
}
