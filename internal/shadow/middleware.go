package shadow

import (
	"io"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps an existing notify pipeline step and additionally
// forwards every outgoing message to a shadow notifier.
type Middleware struct {
	shadow notify.Notifier
	log    io.Writer
}

// NewMiddleware returns a Middleware that shadows all sent messages to
// the provided notifier. log may be nil (defaults to stderr).
func NewMiddleware(shadow notify.Notifier, log io.Writer) *Middleware {
	if log == nil {
		log = defaultWriter()
	}
	return &Middleware{shadow: shadow, log: log}
}

// Apply passes changes through unchanged but fires the shadow notifier
// for each change, logging any errors without blocking.
func (m *Middleware) Apply(changes []state.Change, next func([]state.Change) error) error {
	if m.shadow != nil {
		for _, c := range changes {
			cc := c
			go func() {
				if err := m.shadow.Send(notify.Message{Change: cc}); err != nil {
					logShadowError(m.log, err)
				}
			}()
		}
	}
	return next(changes)
}
