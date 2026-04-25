// Package deadletter provides a dead-letter queue for failed notification deliveries.
// Messages that cannot be delivered after all retry attempts are captured here
// for later inspection or replay.
package deadletter

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Entry holds a failed notification along with metadata.
type Entry struct {
	At      time.Time      `json:"at"`
	Message notify.Message `json:"message"`
	Reason  string         `json:"reason"`
}

// Queue stores undeliverable messages.
type Queue struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
	w       io.Writer
}

// New creates a Queue with the given capacity. If w is nil, os.Stderr is used.
func New(maxSize int, w io.Writer) *Queue {
	if w == nil {
		w = os.Stderr
	}
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Queue{maxSize: maxSize, w: w}
}

// Record adds a failed message to the queue. If the queue is full the oldest
// entry is evicted to make room.
func (q *Queue) Record(msg notify.Message, reason string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	e := Entry{At: time.Now(), Message: msg, Reason: reason}
	if len(q.entries) >= q.maxSize {
		q.entries = q.entries[1:]
	}
	q.entries = append(q.entries, e)

	enc := json.NewEncoder(q.w)
	_ = enc.Encode(e)
}

// Drain returns all queued entries and clears the queue.
func (q *Queue) Drain() []Entry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]Entry, len(q.entries))
	copy(out, q.entries)
	q.entries = q.entries[:0]
	return out
}

// Len returns the current number of queued entries.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.entries)
}
