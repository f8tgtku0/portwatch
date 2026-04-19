// Package rollup groups rapid successive port changes into a single
// summarised notification, reducing alert noise during bulk restarts.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Handler is called with the accumulated batch of changes.
type Handler func(changes []state.Change)

// Rollup collects changes within a sliding window and flushes them together.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	batch   []state.Change
	timer   *time.Timer
	handle  Handler
}

// New creates a Rollup that waits window duration after the last change
// before invoking handle with the accumulated batch.
func New(window time.Duration, handle Handler) *Rollup {
	return &Rollup{
		window: window,
		handle: handle,
	}
}

// Add appends changes to the current batch and resets the flush timer.
func (r *Rollup) Add(changes []state.Change) {
	if len(changes) == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.batch = append(r.batch, changes...)

	if r.timer != nil {
		r.timer.Reset(r.window)
		return
	}
	r.timer = time.AfterFunc(r.window, r.flush)
}

// Flush delivers any pending changes immediately, regardless of the window.
func (r *Rollup) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	batch := r.batch
	r.batch = nil
	r.mu.Unlock()

	if len(batch) > 0 {
		r.handle(batch)
	}
}

func (r *Rollup) flush() {
	r.mu.Lock()
	batch := r.batch
	r.batch = nil
	r.timer = nil
	r.mu.Unlock()

	if len(batch) > 0 {
		r.handle(batch)
	}
}
