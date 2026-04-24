// Package circuit provides a circuit breaker for wrapping notify.Notifier
// implementations. It tracks consecutive send failures and, once a failure
// threshold is reached, opens the circuit to prevent further delivery
// attempts until a configurable cooldown period has elapsed.
//
// Usage:
//
//	base := notify.NewLogNotifier(os.Stdout)
//	breaker := circuit.New(base, 5, 30*time.Second)
//	// breaker implements notify.Notifier
//
// States:
//
//	Closed   — normal operation, all sends forwarded.
//	Open     — failing; sends are rejected immediately.
//	HalfOpen — cooldown elapsed; one probe send is attempted.
package circuit
