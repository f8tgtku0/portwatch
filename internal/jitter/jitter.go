// Package jitter adds randomised delay to notification sends to prevent
// thundering-herd effects when many alerts fire simultaneously.
package jitter

import (
	"context"
	"math/rand"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Jitter holds configuration for randomised delay.
type Jitter struct {
	max time.Duration
	sleep func(context.Context, time.Duration) error
}

// New returns a Jitter that will delay sends by a random duration in [0, max).
// A zero or negative max disables jitter.
func New(max time.Duration) *Jitter {
	return &Jitter{
		max:   max,
		sleep: sleepWithContext,
	}
}

// Delay blocks for a random duration in [0, j.max) then calls next.
// If ctx is cancelled during the sleep the function returns immediately
// without calling next.
func (j *Jitter) Delay(ctx context.Context, msg notify.Message, next func(notify.Message) error) error {
	if j == nil || j.max <= 0 {
		return next(msg)
	}
	//nolint:gosec // non-cryptographic randomness is fine for jitter
	delay := time.Duration(rand.Int63n(int64(j.max)))
	if err := j.sleep(ctx, delay); err != nil {
		return nil // context cancelled — drop silently
	}
	return next(msg)
}

// sleepWithContext sleeps for d or until ctx is done.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
