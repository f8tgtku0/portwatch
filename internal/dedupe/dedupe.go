// Package dedupe provides change deduplication to suppress identical
// consecutive alerts for the same port and action within a time window.
package dedupe

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Deduper tracks recently seen changes and suppresses duplicates.
type Deduper struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     func() time.Time
}

// New creates a Deduper with the given deduplication window.
func New(window time.Duration) *Deduper {
	return &Deduper{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// IsDuplicate returns true if an identical change was seen within the window.
func (d *Deduper) IsDuplicate(c state.Change) bool {
	key := changeKey(c)
	d.mu.Lock()
	defer d.mu.Unlock()
	if t, ok := d.seen[key]; ok && d.now().Sub(t) < d.window {
		return true
	}
	d.seen[key] = d.now()
	return false
}

// Prune removes expired entries from the seen map.
func (d *Deduper) Prune() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Filter returns only changes that are not duplicates.
func (d *Deduper) Filter(changes []state.Change) []state.Change {
	out := make([]state.Change, 0, len(changes))
	for _, c := range changes {
		if !d.IsDuplicate(c) {
			out = append(out, c)
		}
	}
	return out
}

func changeKey(c state.Change) string {
	action := "opened"
	if c.Type == state.Closed {
		action = "closed"
	}
	return fmt.Sprintf("%d:%s", c.Port, action)
}
