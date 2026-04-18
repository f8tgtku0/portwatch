// Package debounce provides a middleware that collapses rapid port-state
// changes (flaps) into a single notification by waiting for a configurable
// quiet window before forwarding events downstream.
//
// Usage:
//
//	mw := debounce.NewMiddleware(500*time.Millisecond, func(changes []state.Change) {
//		// handle stable changes
//	})
//	mw.Apply(changes)
//
// Call Flush() on shutdown to ensure no buffered events are lost.
package debounce
