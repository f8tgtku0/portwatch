package sampling

import (
	"github.com/user/portwatch/internal/state"
)

// Middleware wraps a Sampler and applies it as a pipeline stage.
type Middleware struct {
	sampler *Sampler
}

// NewMiddleware returns a Middleware backed by the given Sampler.
// If sampler is nil, all changes pass through unchanged.
func NewMiddleware(sampler *Sampler) *Middleware {
	return &Middleware{sampler: sampler}
}

// Apply filters changes through the sampler and returns the surviving subset.
// A nil sampler passes all changes through.
func (m *Middleware) Apply(changes []state.Change) []state.Change {
	if m.sampler == nil {
		return changes
	}
	return m.sampler.Sample(changes)
}
