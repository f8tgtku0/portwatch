package suppress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// NotifierMiddleware wraps a notify.Notifier and skips Send calls
// when the current time falls within a suppression window.
type NotifierMiddleware struct {
	inner      notify.Notifier
	suppressor *Suppressor
	log        io.Writer
}

// Wrap returns a NotifierMiddleware around inner using the given Suppressor.
func Wrap(inner notify.Notifier, s *Suppressor) *NotifierMiddleware {
	return &NotifierMiddleware{inner: inner, suppressor: s, log: os.Stderr}
}

// WithWriter sets the writer used for suppression log messages.
func (m *NotifierMiddleware) WithWriter(w io.Writer) *NotifierMiddleware {
	m.log = w
	return m
}

// Send forwards the notification unless suppressed.
func (m *NotifierMiddleware) Send(change state.Change) error {
	if ok, w := m.suppressor.IsSuppressed(time.Now()); ok {
		fmt.Fprintf(m.log, "[suppress] alert silenced (window: %s)\n", w.Reason)
		return nil
	}
	return m.inner.Send(change)
}
