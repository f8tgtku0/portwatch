package audit

import (
	"fmt"

	"github.com/user/portwatch/internal/state"
)

// ChangeAuditor wraps a slice of state.Change values and records each to the Logger.
type ChangeAuditor struct {
	logger *Logger
}

// NewChangeAuditor returns a ChangeAuditor backed by logger.
func NewChangeAuditor(l *Logger) *ChangeAuditor {
	return &ChangeAuditor{logger: l}
}

// Audit records each change as an audit entry and returns the unchanged slice.
func (a *ChangeAuditor) Audit(changes []state.Change) []state.Change {
	for _, c := range changes {
		event := "port.opened"
		if !c.Opened {
			event = "port.closed"
		}
		_ = a.logger.Record(event, c.Port, fmt.Sprintf("proto=%s", c.Proto))
	}
	return changes
}
