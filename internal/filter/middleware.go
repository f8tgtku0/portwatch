package filter

import "github.com/user/portwatch/internal/state"

// Middleware wraps a slice of state.Change values and filters out any
// changes whose port matches the exclusion filter.
type Middleware struct {
	f *Filter
}

// NewMiddleware returns a Middleware that drops changes matching f.
func NewMiddleware(f *Filter) *Middleware {
	return &Middleware{f: f}
}

// Apply returns only the changes whose port is NOT excluded by the filter.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.f == nil || len(changes) == 0 {
		return changes
	}
	out := make([]state.Change, 0, len(changes))
	for _, c := range changes {
		if !m.f.Matches(c.Port) {
			out = append(out, c)
		}
	}
	return out
}
