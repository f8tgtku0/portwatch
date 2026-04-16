package monitor

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ChangeType describes whether a port was opened or closed.
type ChangeType int

const (
	ChangeTypeOpened ChangeType = iota
	ChangeTypeClosed
)

// Change represents a single port state transition.
type Change struct {
	Port int
	Type ChangeType
}

// Monitor periodically scans ports and reports changes.
type Monitor struct {
	scanner  *scanner.Scanner
	interval time.Duration
	previous map[int]bool
}

// New creates a Monitor using the provided scanner and poll interval.
func New(s *scanner.Scanner, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		interval: interval,
		previous: make(map[int]bool),
	}
}

// Diff scans the given port range and returns any changes since the last call.
func (m *Monitor) Diff(start, end int) ([]Change, error) {
	current, err := m.scanner.Scan(start, end)
	if err != nil {
		return nil, err
	}

	currentSet := make(map[int]bool, len(current))
	for _, p := range current {
		currentSet[p] = true
	}

	var changes []Change

	for _, p := range current {
		if !m.previous[p] {
			changes = append(changes, Change{Port: p, Type: ChangeTypeOpened})
		}
	}

	for p := range m.previous {
		if !currentSet[p] {
			changes = append(changes, Change{Port: p, Type: ChangeTypeClosed})
		}
	}

	m.previous = currentSet
	return changes, nil
}

// Reset clears the monitor's previous state, causing the next Diff call to
// treat all open ports as newly opened.
func (m *Monitor) Reset() {
	m.previous = make(map[int]bool)
}

// Run starts the monitoring loop, sending changes to the returned channel.
// The loop stops when done is closed.
func (m *Monitor) Run(start, end int, done <-chan struct{}) <-chan []Change {
	ch := make(chan []Change)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if changes, err := m.Diff(start, end); err == nil && len(changes) > 0 {
					ch <- changes
				}
			}
		}
	}()
	return ch
}
