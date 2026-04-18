// Package debounce delays forwarding of change notifications until
// a quiet period has elapsed, collapsing rapid port flaps into a
// single event.
package debounce

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Debouncer holds pending changes and flushes them after a quiet window.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	pending map[int]state.Change
	timers  map[int]*time.Timer
	flush   func([]state.Change)
}

// New creates a Debouncer that waits window before forwarding changes.
func New(window time.Duration, flush func([]state.Change)) *Debouncer {
	return &Debouncer{
		window:  window,
		pending: make(map[int]state.Change),
		timers:  make(map[int]*time.Timer),
		flush:   flush,
	}
}

// Submit queues a change. If a timer is already running for the port it is
// reset, collapsing the flap.
func (d *Debouncer) Submit(c state.Change) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.pending[c.Port] = c

	if t, ok := d.timers[c.Port]; ok {
		t.Stop()
	}

	d.timers[c.Port] = time.AfterFunc(d.window, func() {
		d.mu.Lock()
		ch, ok := d.pending[c.Port]
		if ok {
			delete(d.pending, c.Port)
			delete(d.timers, c.Port)
		}
		d.mu.Unlock()
		if ok {
			d.flush([]state.Change{ch})
		}
	})
}

// Flush immediately emits all pending changes and cancels outstanding timers.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	changes := make([]state.Change, 0, len(d.pending))
	for _, c := range d.pending {
		changes = append(changes, c)
	}
	for _, t := range d.timers {
		t.Stop()
	}
	d.pending = make(map[int]state.Change)
	d.timers = make(map[int]*time.Timer)
	d.mu.Unlock()

	if len(changes) > 0 {
		d.flush(changes)
	}
}
