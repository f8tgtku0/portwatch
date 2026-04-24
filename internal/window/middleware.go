package window

import (
	"io"
	"log"
	"os"

	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a downstream handler and logs a warning when the number of
// port changes within the sliding window exceeds a configured threshold.
type Middleware struct {
	counter   *Counter
	threshold int
	log       *log.Logger
}

// NewMiddleware creates a Middleware that warns when more than threshold changes
// are observed within the counter's window. If w is nil it defaults to stderr.
func NewMiddleware(counter *Counter, threshold int, w io.Writer) *Middleware {
	if w == nil {
		w = os.Stderr
	}
	return &Middleware{
		counter:   counter,
		threshold: threshold,
		log:       log.New(w, "[window] ", log.LstdFlags),
	}
}

// Apply records each change event and emits a warning if the burst threshold
// is exceeded. It always returns the original changes unmodified.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.counter == nil || len(changes) == 0 {
		return changes
	}
	for range changes {
		m.counter.Record()
	}
	if m.counter.Exceeds(m.threshold) {
		m.log.Printf("burst detected: %d events in window (threshold %d)",
			m.counter.Count(), m.threshold)
	}
	return changes
}
