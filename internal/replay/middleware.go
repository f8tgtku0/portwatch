package replay

import (
	"github.com/user/portwatch/internal/notify"
)

// Middleware wraps a Notifier so that every successfully sent message is
// also recorded in the Replayer buffer for later replay.
type Middleware struct {
	next     notify.Notifier
	replayer *Replayer
}

// NewMiddleware returns a Middleware that forwards to next and records each
// message in r. If r is nil the middleware is a transparent pass-through.
func NewMiddleware(next notify.Notifier, r *Replayer) *Middleware {
	return &Middleware{next: next, replayer: r}
}

// Send forwards the message to the underlying notifier. On success the
// message is recorded in the replay buffer.
func (m *Middleware) Send(msg notify.Message) error {
	if err := m.next.Send(msg); err != nil {
		return err
	}
	if m.replayer != nil {
		m.replayer.Record(msg)
	}
	return nil
}
