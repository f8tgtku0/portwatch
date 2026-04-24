// Package circuit implements a circuit breaker for notification delivery.
// When a notifier fails repeatedly, the circuit opens and stops forwarding
// messages until a cooldown period has elapsed.
package circuit

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota // normal operation
	StateOpen                // failing, blocking sends
	StateHalfOpen            // testing recovery
)

// Breaker is a circuit breaker wrapping a Notifier.
type Breaker struct {
	mu           sync.Mutex
	inner        notify.Notifier
	state        State
	failures     int
	threshold    int
	cooldown     time.Duration
	lastFailure  time.Time
	successes    int
	probeSuccess int
}

// New returns a Breaker that opens after threshold consecutive failures
// and attempts recovery after cooldown.
func New(inner notify.Notifier, threshold int, cooldown time.Duration) *Breaker {
	if threshold <= 0 {
		threshold = 3
	}
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	return &Breaker{
		inner:        inner,
		threshold:    threshold,
		cooldown:     cooldown,
		probeSuccess: 1,
	}
}

// Send forwards the message if the circuit is closed or half-open.
func (b *Breaker) Send(msg notify.Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		if time.Since(b.lastFailure) < b.cooldown {
			return fmt.Errorf("circuit open: too many consecutive failures")
		}
		b.state = StateHalfOpen
		b.successes = 0
	}

	err := b.inner.Send(msg)
	if err != nil {
		b.failures++
		b.lastFailure = time.Now()
		if b.failures >= b.threshold {
			b.state = StateOpen
		}
		return err
	}

	b.successes++
	if b.state == StateHalfOpen && b.successes >= b.probeSuccess {
		b.state = StateClosed
		b.failures = 0
	}
	return nil
}

// CurrentState returns the current circuit state.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Reset forces the circuit back to closed.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = StateClosed
	b.failures = 0
	b.successes = 0
}
