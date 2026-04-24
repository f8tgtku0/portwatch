package quota

import (
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a Quota and filters state.Change slices so that only
// changes whose channel key is within quota are forwarded.
type Middleware struct {
	quota *Quota
	keyer func(c state.Change) string
}

// NewMiddleware returns a Middleware using q to gate changes.
// keyer maps a change to the quota key (typically the notifier channel name).
// If q is nil all changes pass through unchanged.
func NewMiddleware(q *Quota, keyer func(c state.Change) string) *Middleware {
	return &Middleware{quota: q, keyer: keyer}
}

// Apply filters changes, dropping any whose quota key has been exhausted.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.quota == nil || len(changes) == 0 {
		return changes
	}
	out := changes[:0:0]
	for _, c := range changes {
		key := m.keyer(c)
		if m.quota.Allow(key) {
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
