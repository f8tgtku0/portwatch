// Package notify provides notifier implementations for alerting on port changes.
package notify

import "github.com/user/portwatch/internal/state"

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(change state.Change) error
}
