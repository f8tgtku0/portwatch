package grouping

import (
	"github.com/user/portwatch/internal/state"
)

// Middleware attaches group labels to changes via the Labels field.
type Middleware struct {
	grouper *Grouper
}

// NewMiddleware creates a Middleware wrapping the given Grouper.
func NewMiddleware(g *Grouper) *Middleware {
	return &Middleware{grouper: g}
}

// Apply annotates each change with its group label and returns all changes unchanged.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.grouper == nil || len(changes) == 0 {
		return changes
	}
	annotations := m.grouper.Annotate(changes)
	out := make([]state.Change, len(changes))
	for i, c := range changes {
		if label, ok := annotations[c.Port]; ok {
			c.Label = label
		}
		out[i] = c
	}
	return out
}
