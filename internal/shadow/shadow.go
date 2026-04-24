// Package shadow provides a shadow-mode middleware that runs a second
// notifier pipeline in parallel without affecting the primary pipeline.
// This is useful for testing new notification channels in production
// traffic before enabling them fully.
package shadow

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Notifier wraps a primary and a shadow notifier. The shadow notifier
// receives every message but its errors are only logged, never returned.
type Notifier struct {
	primary notify.Notifier
	shadow  notify.Notifier
	log     io.Writer
}

// New returns a Notifier that forwards every message to both primary and
// shadow. Errors from shadow are written to log (defaults to stderr).
func New(primary, shadow notify.Notifier, log io.Writer) *Notifier {
	if log == nil {
		log = os.Stderr
	}
	return &Notifier{primary: primary, shadow: shadow, log: log}
}

// Send delivers the message to the primary notifier and, in a goroutine,
// to the shadow notifier. Only the primary error is returned.
func (n *Notifier) Send(msg notify.Message) error {
	go func() {
		if err := n.shadow.Send(msg); err != nil {
			fmt.Fprintf(n.log, "[shadow] %s notifier error: %v\n",
				time.Now().Format(time.RFC3339), err)
		}
	}()
	return n.primary.Send(msg)
}
