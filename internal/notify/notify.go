// Package notify provides pluggable notification backends for portwatch alerts.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Message holds the data for a single notification event.
type Message struct {
	Level     Level
	Title     string
	Body      string
	Timestamp time.Time
}

// Notifier is the interface that wraps the Send method.
type Notifier interface {
	Send(msg Message) error
}

// LogNotifier writes notifications as structured lines to an io.Writer.
type LogNotifier struct {
	out io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LogNotifier{out: w}
}

// Send formats msg and writes it to the underlying writer.
func (l *LogNotifier) Send(msg Message) error {
	ts := msg.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	_, err := fmt.Fprintf(
		l.out,
		"[%s] %s | %s: %s\n",
		ts.Format(time.RFC3339),
		msg.Level,
		msg.Title,
		msg.Body,
	)
	return err
}

// Multi fans a single message out to multiple Notifiers.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi that dispatches to all provided notifiers.
func NewMulti(notifiers ...Notifier) *Multi {
	return &Multi{notifiers: notifiers}
}

// Send delivers msg to every registered notifier, collecting errors.
func (m *Multi) Send(msg Message) error {
	var first error
	for _, n := range m.notifiers {
		if err := n.Send(msg); err != nil && first == nil {
			first = err
		}
	}
	return first
}
