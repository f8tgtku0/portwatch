// Package enrichment provides port-to-service name resolution for portwatch.
//
// It ships with a built-in registry of well-known TCP/UDP ports and supports
// custom overrides via a user-supplied map. The Enricher is used to annotate
// change events with human-readable service labels before they are passed to
// notifiers or written to the audit log.
//
// Example:
//
//	e := enrichment.New(nil)
//	fmt.Println(e.Label(443)) // "443/https"
package enrichment
