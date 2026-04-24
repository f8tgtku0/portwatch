package fingerprint

import (
	"github.com/user/portwatch/internal/state"
)

// SkipUnchangedMiddleware wraps a downstream Apply func and short-circuits
// processing when the port set has not changed since the last scan cycle.
type SkipUnchangedMiddleware struct {
	fp   *Fingerprinter
	last string
}

// NewMiddleware returns a SkipUnchangedMiddleware backed by fp.
func NewMiddleware(fp *Fingerprinter) *SkipUnchangedMiddleware {
	if fp == nil {
		fp = New()
	}
	return &SkipUnchangedMiddleware{fp: fp}
}

// Apply returns current unchanged when the fingerprint matches the previous
// scan, signalling callers that no diff pipeline run is necessary.
// On a fingerprint mismatch the new hash is stored and current is returned
// as-is so downstream stages can compute the diff.
func (m *SkipUnchangedMiddleware) Apply(current []state.Port) ([]state.Port, bool) {
	h := m.fp.Compute(current)
	if h == m.last {
		return current, false // unchanged
	}
	m.last = h
	return current, true // changed
}

// Reset clears the stored fingerprint, forcing the next call to Apply to
// treat the scan result as a change regardless of content.
func (m *SkipUnchangedMiddleware) Reset() {
	m.last = ""
}
