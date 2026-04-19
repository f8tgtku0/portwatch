// Package redact provides utilities for masking sensitive port ranges
// from alert output and logs (e.g. internal management ports).
package redact

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/state"
)

// Redactor masks changes on ports that are marked sensitive.
type Redactor struct {
	mu     sync.RWMutex
	masked map[int]struct{}
}

// New returns a Redactor that will suppress the given ports from change sets.
func New(ports []int) *Redactor {
	r := &Redactor{masked: make(map[int]struct{}, len(ports))}
	for _, p := range ports {
		r.masked[p] = struct{}{}
	}
	return r
}

// Add marks an additional port as sensitive at runtime.
func (r *Redactor) Add(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.masked[port] = struct{}{}
}

// Remove unmarks a port so changes are visible again.
func (r *Redactor) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.masked, port)
}

// IsMasked reports whether port is currently redacted.
func (r *Redactor) IsMasked(port int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.masked[port]
	return ok
}

// Apply returns a filtered copy of changes with masked ports removed.
// Removed entries are replaced with a placeholder so callers can log counts.
func (r *Redactor) Apply(changes []state.Change) (visible []state.Change, redactedCount int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range changes {
		if _, ok := r.masked[c.Port]; ok {
			redactedCount++
			continue
		}
		visible = append(visible, c)
	}
	return visible, redactedCount
}

// Summary returns a human-readable line describing how many changes were hidden.
func Summary(count int) string {
	if count == 0 {
		return ""
	}
	return fmt.Sprintf("[redact] %d change(s) hidden (sensitive ports)", count)
}
