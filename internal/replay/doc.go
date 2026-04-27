// Package replay provides a buffer-backed mechanism for recording and
// re-delivering notify.Message values that have already been sent through
// the notification pipeline.
//
// Use cases include:
//   - Replaying missed alerts after a notifier outage.
//   - Smoke-testing a newly configured notification channel against
//     real historical events.
//   - Auditing what was sent during a maintenance window.
//
// The Replayer is safe for concurrent use. The optional Middleware wrapper
// ensures every successfully delivered message is automatically buffered.
package replay
