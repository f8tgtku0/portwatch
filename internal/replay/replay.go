// Package replay provides a mechanism to re-deliver previously recorded
// change events to a notifier, useful for catching up missed alerts or
// re-testing notification pipelines.
package replay

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Entry holds a single recorded change message with its timestamp.
type Entry struct {
	At      time.Time
	Message notify.Message
}

// Replayer stores entries and can replay them to a notifier.
type Replayer struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
	w       io.Writer
}

// New creates a Replayer with the given maximum buffer size.
// If w is nil, os.Stderr is used for diagnostic output.
func New(maxSize int, w io.Writer) *Replayer {
	if w == nil {
		w = os.Stderr
	}
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Replayer{maxSize: maxSize, w: w}
}

// Record appends a message to the replay buffer, evicting the oldest entry
// when the buffer is full.
func (r *Replayer) Record(msg notify.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) >= r.maxSize {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, Entry{At: time.Now(), Message: msg})
}

// Replay delivers all buffered entries to n in chronological order.
// Errors from the notifier are logged but do not halt delivery.
func (r *Replayer) Replay(n notify.Notifier) int {
	r.mu.Lock()
	copy := make([]Entry, len(r.entries))
	copy_ := r.entries
	_ = copy
	entries := make([]Entry, len(copy_))
	copy(entries, copy_)
	r.mu.Unlock()

	delivered := 0
	for _, e := range entries {
		if err := n.Send(e.Message); err != nil {
			fmt.Fprintf(r.w, "replay: send error for port %d: %v\n", e.Message.Port, err)
			continue
		}
		delivered++
	}
	return delivered
}

// Len returns the number of buffered entries.
func (r *Replayer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

// Clear removes all buffered entries.
func (r *Replayer) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = r.entries[:0]
}
