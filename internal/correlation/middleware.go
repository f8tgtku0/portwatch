package correlation

import (
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a Correlator and implements a pipeline Apply step.
// It feeds each batch of changes into the correlator and immediately
// drains any ready events, returning their flattened changes.
type Middleware struct {
	c *Correlator
}

// NewMiddleware creates a Middleware backed by the given Correlator.
// Pass nil to disable correlation (all changes pass through unchanged).
func NewMiddleware(c *Correlator) *Middleware {
	return &Middleware{c: c}
}

// Apply feeds changes into the correlator and returns any events that
// have already fired (non-blocking drain). If the correlator is nil,
// all changes are returned as-is.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.c == nil {
		return changes
	}

	for _, ch := range changes {
		m.c.Add(ch)
	}

	var out []state.Change
	for {
		select {
		case ev := <-m.c.Events():
			out = append(out, ev.Changes...)
		default:
			return out
		}
	}
}
