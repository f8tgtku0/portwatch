package baseline

import (
	"github.com/user/portwatch/internal/state"
)

// Middleware filters out changes whose ports are present in the baseline.
type Middleware struct {
	b *Baseline
}

// NewMiddleware returns a Middleware wrapping the given Baseline.
func NewMiddleware(b *Baseline) *Middleware {
	return &Middleware{b: b}
}

// Apply removes changes that involve ports present in the baseline.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.b == nil {
		return changes
	}
	out := changes[:0:0]
	for _, c := range changes {
		if !m.b.Contains(c.Port) {
			out = append(out, c)
		}
	}
	return out
}
