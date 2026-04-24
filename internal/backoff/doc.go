// Package backoff implements exponential backoff with optional jitter.
//
// It is used by portwatch to gracefully handle transient failures when
// delivering alerts to notification channels (webhooks, email, SMS, etc.).
//
// Usage:
//
//	p := backoff.DefaultPolicy()
//	err := p.Runner(3, time.Sleep, func() error {
//		return notifier.Send(change)
//	})
package backoff
