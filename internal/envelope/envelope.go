// Package envelope wraps outgoing notify messages with metadata such as
// hostname, environment tag, and a monotonic send sequence number so that
// receivers can correlate and deduplicate alerts across multiple portwatch
// instances running in the same environment.
package envelope

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Envelope wraps a notify.Message with additional routing metadata.
type Envelope struct {
	// Message is the original change notification.
	Message notify.Message

	// Host is the hostname of the machine that generated the alert.
	Host string

	// Env is an optional environment label (e.g. "prod", "staging").
	Env string

	// Seq is a per-process monotonically increasing sequence number.
	Seq uint64

	// WrappedAt is the wall-clock time at which the envelope was created.
	WrappedAt time.Time
}

// String returns a human-readable summary of the envelope metadata.
func (e Envelope) String() string {
	env := e.Env
	if env == "" {
		env = "(none)"
	}
	return fmt.Sprintf("[seq=%d host=%s env=%s at=%s]",
		e.Seq, e.Host, env, e.WrappedAt.Format(time.RFC3339))
}

// Wrapper is a notify.Notifier middleware that wraps each message in an
// Envelope before forwarding it to an inner notifier that understands the
// enriched format via its Subject/Body fields.
type Wrapper struct {
	inner notify.Notifier
	host  string
	env   string
	seq   atomic.Uint64
}

// New returns a Wrapper that decorates inner with envelope metadata.
// env may be empty; host defaults to the OS hostname when empty.
func New(inner notify.Notifier, env string) (*Wrapper, error) {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return &Wrapper{
		inner: inner,
		host:  host,
		env:   env,
	}, nil
}

// Send wraps msg in an Envelope, rewrites its Subject and Body to include
// the metadata prefix, and forwards the enriched message to the inner notifier.
func (w *Wrapper) Send(msg notify.Message) error {
	if w.inner == nil {
		return nil
	}

	env := w.wrap(msg)

	enriched := notify.Message{
		Port:      msg.Port,
		Action:    msg.Action,
		Timestamp: msg.Timestamp,
		Label:     msg.Label,
		Subject:   fmt.Sprintf("%s %s", env.String(), msg.Subject),
		Body:      fmt.Sprintf("%s\n%s", env.String(), msg.Body),
	}

	return w.inner.Send(enriched)
}

// wrap builds an Envelope for msg and increments the internal counter.
func (w *Wrapper) wrap(msg notify.Message) Envelope {
	return Envelope{
		Message:   msg,
		Host:      w.host,
		Env:       w.env,
		Seq:       w.seq.Add(1),
		WrappedAt: time.Now(),
	}
}
