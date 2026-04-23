// Package retry wraps a notify.Notifier with automatic retry logic.
//
// Usage:
//
//	policy := retry.DefaultPolicy()
//	policy.MaxAttempts = 5
//	policy.Delay = time.Second
//
//	rn := retry.New(baseNotifier, policy)
//	// rn.Send will retry up to 5 times before returning an error.
package retry
