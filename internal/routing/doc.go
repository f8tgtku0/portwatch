// Package routing maps observed open ports to known network service routing
// hints, including direction (inbound/outbound), protocol, and a human-readable
// service note.
//
// Built-in mappings cover common well-known ports (SSH, HTTP, HTTPS, MySQL,
// PostgreSQL, Redis, etc.). User-defined overrides take priority over built-ins,
// allowing custom annotations for non-standard or internal services.
//
// Usage:
//
//	r := routing.New(overrides)
//	rt := r.Lookup(443)  // Route{Note: "HTTPS", ...}
//	annotated := r.Annotate(ports)
package routing
