package throttle

import (
	"context"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a Notifier and suppresses sends that exceed the throttle.
type Middleware struct {
	next     notify.Notifier
	throttle *Throttle
}

// NewMiddleware returns a throttling middleware around the given Notifier.
func NewMiddleware(n notify.Notifier, th *Throttle) *Middleware {
	return &Middleware{next: n, throttle: th}
}

// Send forwards the message only if the throttle allows it.
func (m *Middleware) Send(ctx context.Context, msg notify.Message) error {
	key := changeKey(msg.Change)
	if !m.throttle.Allow(key) {
		return nil
	}
	return m.next.Send(ctx, msg)
}

func changeKey(c state.Change) string {
	action := "opened"
	if !c.Opened {
		action = "closed"
	}
	return fmt.Sprintf("port:%d:%s", c.Port, action)
}
