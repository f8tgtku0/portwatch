// Package triage classifies port changes by severity based on
// configurable rules (port ranges, known services, direction).
package triage

import (
	"github.com/user/portwatch/internal/state"
)

// Level represents an alert severity level.
type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelWarning:
		return "warning"
	case LevelCritical:
		return "critical"
	default:
		return "info"
	}
}

// Rule maps a port range and action to a severity level.
type Rule struct {
	MinPort int
	MaxPort int
	Opened  *Level
	Closed  *Level
}

// Triage classifies changes according to a set of rules.
type Triage struct {
	rules []Rule
}

// New creates a Triage with the given rules.
func New(rules []Rule) *Triage {
	return &Triage{rules: rules}
}

// Classify returns the severity level for a given change.
func (t *Triage) Classify(c state.Change) Level {
	for _, r := range t.rules {
		if c.Port < r.MinPort || c.Port > r.MaxPort {
			continue
		}
		if c.Action == state.Opened && r.Opened != nil {
			return *r.Opened
		}
		if c.Action == state.Closed && r.Closed != nil {
			return *r.Closed
		}
	}
	return LevelInfo
}

// Annotate returns a map of change to its classified level.
func (t *Triage) Annotate(changes []state.Change) map[state.Change]Level {
	out := make(map[state.Change]Level, len(changes))
	for _, c := range changes {
		out[c] = t.Classify(c)
	}
	return out
}
