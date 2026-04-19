package enrichment

import "github.com/user/portwatch/internal/state"

// Middleware attaches service name labels to port change messages.
type Middleware struct {
	lookup *Lookup
}

// NewMiddleware returns a new enrichment Middleware using the given Lookup.
func NewMiddleware(l *Lookup) *Middleware {
	return &Middleware{lookup: l}
}

// Enrich annotates each change in the slice with a human-readable service
// label derived from the port number and returns the annotated slice.
// The original slice is not modified; a new slice is returned.
func (m *Middleware) Enrich(changes []state.Change) []EnrichedChange {
	out := make([]EnrichedChange, 0, len(changes))
	for _, c := range changes {
		label := ""
		if m.lookup != nil {
			label = m.lookup.Label(c.Port)
		}
		out = append(out, EnrichedChange{
			Change:  c,
			Service: label,
		})
	}
	return out
}

// EnrichedChange wraps a state.Change with an optional service label.
type EnrichedChange struct {
	state.Change
	Service string
}
