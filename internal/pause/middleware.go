package pause

import (
	"context"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a Notifier and silently drops messages while paused.
type Middleware struct {
	pauser *Pauser
	next   notify.Notifier
}

// NewMiddleware returns a Middleware that gates delivery through p.
// When p.IsPaused() is true, Send returns nil without forwarding.
func NewMiddleware(p *Pauser, next notify.Notifier) *Middleware {
	if p == nil {
		return &Middleware{pauser: New(), next: next}
	}
	return &Middleware{pauser: p, next: next}
}

// Send forwards the message to the underlying notifier unless paused.
func (m *Middleware) Send(ctx context.Context, msg notify.Message) error {
	if m.pauser.IsPaused() {
		return nil
	}
	return m.next.Send(ctx, msg)
}

// Apply filters a slice of state.Change values, returning nil when paused.
// This satisfies the optional pipeline Apply convention used by other middleware.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.pauser.IsPaused() {
		return nil
	}
	return changes
}
