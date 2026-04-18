package debounce

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a downstream flush function with debounce logic.
type Middleware struct {
	debouncer *Debouncer
}

// NewMiddleware returns a Middleware that debounces changes by window before
// passing them to next.
func NewMiddleware(window time.Duration, next func([]state.Change)) *Middleware {
	return &Middleware{
		debouncer: New(window, next),
	}
}

// Apply submits each change to the debouncer.
func (m *Middleware) Apply(changes []state.Change) {
	for _, c := range changes {
		m.debouncer.Submit(c)
	}
}

// Flush forces immediate delivery of all buffered changes.
func (m *Middleware) Flush() {
	m.debouncer.Flush()
}
