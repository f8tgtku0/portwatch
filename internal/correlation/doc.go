// Package correlation groups rapid successive changes on the same port
// into a single correlated Event. This prevents alert storms when a
// service bounces — closing and reopening a port within a short window
// is emitted as one logical event rather than two separate alerts.
//
// Usage:
//
//	c := correlation.New(2 * time.Second)
//	mw := correlation.NewMiddleware(c)
//	resolved := mw.Apply(changes)
package correlation
