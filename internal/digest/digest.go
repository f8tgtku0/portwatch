// Package digest provides periodic summary digests of port change activity.
package digest

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Entry records a single change event for digest aggregation.
type Entry struct {
	Port   int
	Action string // "opened" or "closed"
	At     time.Time
}

// Digest accumulates change entries and emits periodic summaries.
type Digest struct {
	mu      sync.Mutex
	entries []Entry
	writer  io.Writer
}

// New creates a new Digest writing summaries to w.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Digest {
	if w == nil {
		w = os.Stdout
	}
	return &Digest{writer: w}
}

// Record adds change entries derived from a state diff result.
func (d *Digest) Record(opened, closed []state.Port) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now()
	for _, p := range opened {
		d.entries = append(d.entries, Entry{Port: p.Number, Action: "opened", At: now})
	}
	for _, p := range closed {
		d.entries = append(d.entries, Entry{Port: p.Number, Action: "closed", At: now})
	}
}

// Flush writes a summary of accumulated entries to the writer and clears them.
func (d *Digest) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.entries) == 0 {
		return
	}
	opened, closed := 0, 0
	for _, e := range d.entries {
		if e.Action == "opened" {
			opened++
		} else {
			closed++
		}
	}
	fmt.Fprintf(d.writer, "[digest] %s — %d opened, %d closed (%d total changes)\n",
		time.Now().Format(time.RFC3339), opened, closed, len(d.entries))
	for _, e := range d.entries {
		fmt.Fprintf(d.writer, "  port %d %s at %s\n", e.Port, e.Action, e.At.Format(time.RFC3339))
	}
	d.entries = d.entries[:0]
}

// Len returns the number of buffered entries.
func (d *Digest) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.entries)
}
