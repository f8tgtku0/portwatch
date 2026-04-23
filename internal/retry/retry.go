// Package retry provides a middleware that retries failed notify sends
// with configurable attempts and backoff.
package retry

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Policy defines how retries are performed.
type Policy struct {
	// MaxAttempts is the total number of send attempts (including the first).
	MaxAttempts int
	// Delay is the wait duration between attempts.
	Delay time.Duration
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
	}
}

// Notifier wraps a notify.Notifier and retries on failure.
type Notifier struct {
	inner  notify.Notifier
	policy Policy
	sleep  func(time.Duration)
}

// New returns a Notifier that retries the inner notifier according to policy.
func New(inner notify.Notifier, policy Policy) *Notifier {
	return &Notifier{
		inner:  inner,
		policy: policy,
		sleep:  time.Sleep,
	}
}

// Send attempts to deliver msg, retrying up to policy.MaxAttempts times.
// It returns the last error if all attempts fail.
func (n *Notifier) Send(msg notify.Message) error {
	if n.inner == nil {
		return nil
	}
	max := n.policy.MaxAttempts
	if max < 1 {
		max = 1
	}
	var err error
	for attempt := 1; attempt <= max; attempt++ {
		err = n.inner.Send(msg)
		if err == nil {
			return nil
		}
		if attempt < max {
			n.sleep(n.policy.Delay)
		}
	}
	return fmt.Errorf("retry: all %d attempts failed: %w", max, err)
}
