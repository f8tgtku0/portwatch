// Package filter implements port exclusion rules for portwatch.
//
// A Filter is constructed from a slice of port specifications, each of which
// is either a single port number (e.g. "22") or an inclusive range
// (e.g. "8000-9000"). Ports matching any rule are silently skipped during
// change detection so that known, expected services do not generate alerts.
//
// Example usage:
//
//	f, err := filter.New(cfg.IgnorePorts)
//	if err != nil { ... }
//	if f.Ignored(port) {
//	    continue
//	}
package filter
