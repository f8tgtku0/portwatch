// Package watchdog provides a self-monitoring component that detects
// when the scan loop stalls or falls behind schedule.
package watchdog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Watchdog monitors the scan loop heartbeat and fires an alert when
// no heartbeat is received within the configured timeout.
type Watchdog struct {
	timeout  time.Duration
	ticker   *time.Ticker
	lastBeat time.Time
	mu       sync.Mutex
	w        io.Writer
	stop     chan struct{}
}

// New creates a Watchdog that fires if no heartbeat is received within timeout.
func New(timeout time.Duration, w io.Writer) *Watchdog {
	if w == nil {
		w = os.Stderr
	}
	return &Watchdog{
		timeout: timeout,
		w:       w,
		stop:    make(chan struct{}),
	}
}

// Beat records a heartbeat from the scan loop.
func (wd *Watchdog) Beat() {
	wd.mu.Lock()
	wd.lastBeat = time.Now()
	wd.mu.Unlock()
}

// Start begins watching for missed heartbeats. It blocks until Stop is called.
func (wd *Watchdog) Start() {
	wd.mu.Lock()
	wd.lastBeat = time.Now()
	wd.mu.Unlock()

	wd.ticker = time.NewTicker(wd.timeout / 2)
	defer wd.ticker.Stop()

	for {
		select {
		case <-wd.stop:
			return
		case t := <-wd.ticker.C:
			wd.mu.Lock()
			last := wd.lastBeat
			wd.mu.Unlock()
			if t.Sub(last) > wd.timeout {
				fmt.Fprintf(wd.w, "[watchdog] WARNING: scan loop stalled — last heartbeat %s ago\n",
					t.Sub(last).Round(time.Second))
			}
		}
	}
}

// Stop halts the watchdog.
func (wd *Watchdog) Stop() {
	close(wd.stop)
}

// LastBeat returns the time of the most recent heartbeat.
func (wd *Watchdog) LastBeat() time.Time {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	return wd.lastBeat
}
