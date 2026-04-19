// Package correlation groups related port changes into correlated events.
// For example, a service restart may close then reopen the same port within
// a short window — correlation surfaces that as a single logical event.
package correlation

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Event holds a set of changes that are considered correlated.
type Event struct {
	ID      string
	Changes []state.Change
	At      time.Time
}

// Correlator buffers changes and groups those affecting the same port
// within a sliding window into a single Event.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	buckets map[int][]state.Change
	timers  map[int]*time.Timer
	out     chan Event
}

// New creates a Correlator with the given window duration.
func New(window time.Duration) *Correlator {
	return &Correlator{
		window:  window,
		buckets: make(map[int][]state.Change),
		timers:  make(map[int]*time.Timer),
		out:     make(chan Event, 64),
	}
}

// Add buffers a change. If no timer exists for the port, one is started;
// additional changes within the window extend the bucket.
func (c *Correlator) Add(ch state.Change) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buckets[ch.Port] = append(c.buckets[ch.Port], ch)

	if t, ok := c.timers[ch.Port]; ok {
		t.Reset(c.window)
		return
	}

	port := ch.Port
	c.timers[port] = time.AfterFunc(c.window, func() {
		c.flush(port)
	})
}

// Events returns the channel on which correlated Events are delivered.
func (c *Correlator) Events() <-chan Event {
	return c.out
}

// Flush immediately emits any buffered changes for all ports.
func (c *Correlator) Flush() {
	c.mu.Lock()
	ports := make([]int, 0, len(c.buckets))
	for p := range c.buckets {
		ports = append(ports, p)
	}
	c.mu.Unlock()

	for _, p := range ports {
		c.flush(p)
	}
}

func (c *Correlator) flush(port int) {
	c.mu.Lock()
	changes := c.buckets[port]
	delete(c.buckets, port)
	if t, ok := c.timers[port]; ok {
		t.Stop()
		delete(c.timers, port)
	}
	c.mu.Unlock()

	if len(changes) == 0 {
		return
	}
	c.out <- Event{
		ID:      fmt.Sprintf("port-%d-%d", port, time.Now().UnixNano()),
		Changes: changes,
		At:      time.Now(),
	}
}
