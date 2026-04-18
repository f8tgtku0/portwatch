// Package audit provides a structured audit log for portwatch events,
// recording who triggered actions and when for compliance purposes.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Port      int       `json:"port,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	mu sync.Mutex
	w  io.Writer
}

// New returns a Logger writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Record writes a single audit entry.
func (l *Logger) Record(event string, port int, detail string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Port:      port,
		Detail:    detail,
	}
	return l.write(e)
}

func (l *Logger) write(e Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
