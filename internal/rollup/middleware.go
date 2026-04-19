package rollup

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a downstream handler with rollup batching.
type Middleware struct {
	rollup *Rollup
}

// NewMiddleware returns a Middleware that buffers changes for window duration
// before forwarding the accumulated batch to next.
func NewMiddleware(window time.Duration, next func([]state.Change)) *Middleware {
	return &Middleware{
		rollup: New(window, next),
	}
}

// Apply adds changes to the rollup buffer.
func (m *Middleware) Apply(changes []state.Change) {
	m.rollup.Add(changes)
}

// Flush forces immediate delivery of any buffered changes.
func (m *Middleware) Flush() {
	m.rollup.Flush()
}
