// Package backoff provides exponential backoff with jitter for retry logic
// when sending notifications or scanning ports under load.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Policy defines the parameters for exponential backoff.
type Policy struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	JitterFraction  float64 // 0.0 = no jitter, 1.0 = full jitter
}

// DefaultPolicy returns a sensible backoff policy suitable for notification retries.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0.3,
	}
}

// Backoff computes the wait duration for a given attempt (0-indexed).
// It applies exponential growth capped at MaxInterval, then adds jitter.
func (p Policy) Backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	base := float64(p.InitialInterval) * math.Pow(p.Multiplier, float64(attempt))
	if base > float64(p.MaxInterval) {
		base = float64(p.MaxInterval)
	}
	jitter := base * p.JitterFraction * (rand.Float64()*2 - 1)
	result := time.Duration(base + jitter)
	if result < 0 {
		result = 0
	}
	return result
}

// Sleeper abstracts time.Sleep to allow testing without real delays.
type Sleeper func(time.Duration)

// Runner executes fn up to maxAttempts times, sleeping between failures.
// It returns the last error if all attempts fail.
func (p Policy) Runner(maxAttempts int, sleep Sleeper, fn func() error) error {
	if maxAttempts <= 0 {
		maxAttempts = 1
	}
	var err error
	for i := 0; i < maxAttempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i < maxAttempts-1 {
			sleep(p.Backoff(i))
		}
	}
	return err
}
