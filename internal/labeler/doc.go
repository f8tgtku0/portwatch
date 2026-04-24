// Package labeler resolves human-readable service names for port numbers.
//
// It combines a built-in table of IANA well-known port assignments with an
// optional user-supplied override map so that alert messages can include
// context such as "ssh" or "postgres" alongside the raw port number.
//
// Usage:
//
//	l := labeler.New(map[int]string{8888: "jupyter"})
//	name := l.Lookup(22)   // "ssh"
//	name  = l.Lookup(8888) // "jupyter"
//
// The Middleware type wraps a notify.Notifier and automatically annotates
// every outbound message before forwarding it downstream.
package labeler
