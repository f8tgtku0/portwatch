package labeler

import (
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a notify.Notifier and enriches each message with a service
// label before forwarding it downstream.
type Middleware struct {
	labeler  *Labeler
	next     notify.Notifier
}

// NewMiddleware returns a Middleware that labels changes using l before
// passing them to next. If l is nil all changes are forwarded as-is.
func NewMiddleware(l *Labeler, next notify.Notifier) *Middleware {
	return &Middleware{labeler: l, next: next}
}

// Apply labels the changes in msg and forwards the enriched message.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.labeler == nil {
		return changes
	}
	return m.labeler.Annotate(changes)
}

// Send enriches msg and forwards it to the wrapped notifier.
func (m *Middleware) Send(msg notify.Message) error {
	if m.labeler != nil {
		msg.Changes = m.labeler.Annotate(msg.Changes)
	}
	return m.next.Send(msg)
}
