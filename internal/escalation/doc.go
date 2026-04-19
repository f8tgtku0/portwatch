// Package escalation implements multi-tier alert escalation for portwatch.
//
// When a port change is detected and remains unacknowledged, the escalator
// will re-notify through progressively more urgent channels after configurable
// time thresholds.
//
// Example usage:
//
//	esc := escalation.New([]escalation.Tier{
//		{After: 5 * time.Minute, Notifier: slackNotifier},
//		{After: 15 * time.Minute, Notifier: pagerdutyNotifier},
//	})
//	esc.Track("opened:8080")
//	// later...
//	esc.Evaluate(change)
//	esc.Acknowledge("opened:8080")
package escalation
