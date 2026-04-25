// Package deadletter implements a bounded dead-letter queue for portwatch
// notifications.
//
// When a downstream notifier fails to deliver a message the middleware wraps
// the error, stores the original message with a timestamp and failure reason,
// and surfaces the entry for later inspection or replay.
//
// Usage:
//
//	q  := deadletter.New(200, os.Stderr)
//	mw := deadletter.NewMiddleware(slackNotifier, q)
//	// mw satisfies notify.Notifier
//
//	// later – drain and re-deliver
//	for _, e := range q.Drain() {
//		_ = slackNotifier.Send(e.Message)
//	}
package deadletter
