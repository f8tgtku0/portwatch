// Package truncate provides a middleware that limits the number of changes
// forwarded to a notifier in a single batch. This protects downstream systems
// from being overwhelmed when a large burst of port changes occurs.
package truncate

import (
	"fmt"
	"io"
	"os"
	"time"

	"portwatch/internal/state"
)

// Truncator limits a slice of changes to at most MaxSize entries.
type Truncator struct {
	maxSize int
	writer  io.Writer
}

// New creates a Truncator that will forward at most maxSize changes per batch.
// If maxSize is <= 0 it is clamped to 1. A nil writer defaults to os.Stderr.
func New(maxSize int, w io.Writer) *Truncator {
	if maxSize <= 0 {
		maxSize = 1
	}
	if w == nil {
		w = os.Stderr
	}
	return &Truncator{maxSize: maxSize, writer: w}
}

// Apply returns up to maxSize changes from the provided slice. If the input
// exceeds the limit a warning is written to the configured writer and the
// excess entries are silently dropped.
func (t *Truncator) Apply(changes []state.Change) []state.Change {
	if len(changes) <= t.maxSize {
		return changes
	}
	dropped := len(changes) - t.maxSize
	fmt.Fprintf(t.writer, "%s [truncate] batch truncated: keeping %d of %d changes (%d dropped)\n",
		time.Now().UTC().Format(time.RFC3339), t.maxSize, len(changes), dropped)
	return changes[:t.maxSize]
}

// MaxSize returns the configured maximum batch size.
func (t *Truncator) MaxSize() int {
	return t.maxSize
}
