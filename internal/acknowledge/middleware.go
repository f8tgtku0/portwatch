package acknowledge

import (
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a Notifier and silently drops any change messages that
// have been acknowledged in the provided Store.
type Middleware struct {
	store *Store
	next  notify.Notifier
}

// NewMiddleware returns a Middleware that filters acknowledged changes before
// forwarding the remainder to next. If store is nil all changes are forwarded.
func NewMiddleware(store *Store, next notify.Notifier) *Middleware {
	return &Middleware{store: store, next: next}
}

// Send filters out acknowledged changes and forwards the rest.
func (m *Middleware) Send(msg notify.Message) error {
	if m.store == nil {
		return m.next.Send(msg)
	}

	filtered := make([]state.Change, 0, len(msg.Changes))
	for _, c := range msg.Changes {
		if !m.store.IsAcknowledged(c.Port, string(c.Action)) {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	msg.Changes = filtered
	return m.next.Send(msg)
}
