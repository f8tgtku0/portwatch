package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event holds information about a port change.
type Event struct {
	Timestamp time.Time
	Level     Level
	Change    monitor.Change
}

// Notifier sends alerts for port change events.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert for the given change.
func (n *Notifier) Notify(c monitor.Change) error {
	ev := Event{
		Timestamp: time.Now(),
		Change:    c,
	}

	switch c.Type {
	case monitor.ChangeTypeOpened:
		ev.Level = LevelAlert
		_, err := fmt.Fprintf(n.out, "[%s] %s — port %d OPENED\n",
			ev.Timestamp.Format(time.RFC3339), ev.Level, c.Port)
		return err
	case monitor.ChangeTypeClosed:
		ev.Level = LevelWarn
		_, err := fmt.Fprintf(n.out, "[%s] %s — port %d CLOSED\n",
			ev.Timestamp.Format(time.RFC3339), ev.Level, c.Port)
		return err
	default:
		return fmt.Errorf("unknown change type: %v", c.Type)
	}
}

// NotifyAll calls Notify for each change in the slice.
func (n *Notifier) NotifyAll(changes []monitor.Change) error {
	for _, c := range changes {
		if err := n.Notify(c); err != nil {
			return err
		}
	}
	return nil
}
