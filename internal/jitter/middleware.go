package jitter

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Middleware wraps a Notifier and injects a random delay before each send.
type Middleware struct {
	jitter *Jitter
	next   notify.Notifier
}

// NewMiddleware returns a Middleware that delays sends by up to max before
// forwarding to next. A nil or zero-max jitter passes messages through
// immediately.
func NewMiddleware(max time.Duration, next notify.Notifier) *Middleware {
	return &Middleware{
		jitter: New(max),
		next:   next,
	}
}

// Send applies jitter delay then forwards msg to the wrapped Notifier.
func (m *Middleware) Send(msg notify.Message) error {
	if m.jitter == nil || m.next == nil {
		if m.next != nil {
			return m.next.Send(msg)
		}
		return nil
	}
	return m.jitter.Delay(context.Background(), msg, m.next.Send)
}
