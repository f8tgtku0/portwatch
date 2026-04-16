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

// levelForChange returns the alert Level appropriate for the given ChangeType.
func levelForChange(t monitor.ChangeType) (Level, error) {
	switch t {
	case monitor.ChangeTypeOpened:
		return LevelAlert, nil
	case monitor.ChangeTypeClosed:
		return LevelWarn, nil
	default:
		return "", fmt.Errorf("unknown change type: %v", t)
	}
}

// Notify formats and writes an alert for the given change.
func (n *Notifier) Notify(c monitor.Change) error {
	level, err := levelForChange(c.Type)
	if err != nil {
		return err
	}

	ev := Event{
		Timestamp: time.Now(),
		Level:     level,
		Change:    c,
	}

	_, err = fmt.Fprintf(n.out, "[%s] %s — port %d %s\n",
		ev.Timestamp.Format(time.RFC3339), ev.Level, c.Port, c.Type)
	return err
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
