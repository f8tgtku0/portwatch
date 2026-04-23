package watchdog

import "github.com/user/portwatch/internal/state"

// BeatMiddleware wraps a scan result handler and records a heartbeat
// each time it is invoked, keeping the Watchdog satisfied.
type BeatMiddleware struct {
	wd   *Watchdog
	next func([]state.PortChange) []state.PortChange
}

// NewBeatMiddleware returns a BeatMiddleware that calls wd.Beat() on every
// invocation before delegating to next. If next is nil the changes are
// returned unmodified.
func NewBeatMiddleware(wd *Watchdog, next func([]state.PortChange) []state.PortChange) *BeatMiddleware {
	return &BeatMiddleware{wd: wd, next: next}
}

// Apply records a heartbeat and forwards changes to the next handler.
func (m *BeatMiddleware) Apply(changes []state.PortChange) []state.PortChange {
	if m.wd != nil {
		m.wd.Beat()
	}
	if m.next != nil {
		return m.next(changes)
	}
	return changes
}
