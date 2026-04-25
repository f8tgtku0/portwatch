package deadletter

import (
	"fmt"

	"github.com/user/portwatch/internal/notify"
)

// Middleware wraps a Notifier and routes failed sends to a dead-letter Queue.
type Middleware struct {
	next  notify.Notifier
	queue *Queue
}

// NewMiddleware returns a Middleware that forwards to next and captures
// any send errors in q.
func NewMiddleware(next notify.Notifier, q *Queue) notify.Notifier {
	if next == nil {
		panic("deadletter: next notifier must not be nil")
	}
	if q == nil {
		panic("deadletter: queue must not be nil")
	}
	return &Middleware{next: next, queue: q}
}

// Send attempts delivery via the wrapped notifier. On failure the message is
// recorded in the dead-letter queue and the original error is returned so the
// caller can still react to it.
func (m *Middleware) Send(msg notify.Message) error {
	if err := m.next.Send(msg); err != nil {
		m.queue.Record(msg, fmt.Sprintf("%v", err))
		return err
	}
	return nil
}
