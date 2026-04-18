// Package audit provides structured audit logging for portwatch.
//
// Each port change event (opened or closed) is recorded as a
// newline-delimited JSON entry containing a UTC timestamp, event name,
// port number, and optional detail string.
//
// Usage:
//
//	f, _ := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	l := audit.New(f)
//	ca := audit.NewChangeAuditor(l)
//	changes = ca.Audit(changes) // records and passes changes through
package audit
