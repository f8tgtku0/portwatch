// Package triage classifies port change events into severity levels
// (info, warning, critical) based on user-defined rules.
//
// Each Rule binds a port range to optional severity levels for opened
// and closed events. The first matching rule wins. Changes that match
// no rule default to LevelInfo.
//
// Example usage:
//
//	critical := triage.LevelCritical
//	warning  := triage.LevelWarning
//	tr := triage.New([]triage.Rule{
//		{MinPort: 1, MaxPort: 1023, Opened: &critical, Closed: &warning},
//	})
//	level := tr.Classify(change)
package triage
