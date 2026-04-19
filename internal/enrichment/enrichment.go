// Package enrichment resolves service names for well-known ports.
package enrichment

import "fmt"

// Entry holds enriched metadata for a port change.
type Entry struct {
	Port    int
	Service string
	Proto   string
}

// Enricher resolves port metadata.
type Enricher struct {
	custom map[int]Entry
}

// wellKnown contains a small built-in registry of common ports.
var wellKnown = map[int]Entry{
	21:   {Port: 21, Service: "ftp", Proto: "tcp"},
	22:   {Port: 22, Service: "ssh", Proto: "tcp"},
	23:   {Port: 23, Service: "telnet", Proto: "tcp"},
	25:   {Port: 25, Service: "smtp", Proto: "tcp"},
	53:   {Port: 53, Service: "dns", Proto: "udp"},
	80:   {Port: 80, Service: "http", Proto: "tcp"},
	110:  {Port: 110, Service: "pop3", Proto: "tcp"},
	143:  {Port: 143, Service: "imap", Proto: "tcp"},
	443:  {Port: 443, Service: "https", Proto: "tcp"},
	3306: {Port: 3306, Service: "mysql", Proto: "tcp"},
	5432: {Port: 5432, Service: "postgres", Proto: "tcp"},
	6379: {Port: 6379, Service: "redis", Proto: "tcp"},
	8080: {Port: 8080, Service: "http-alt", Proto: "tcp"},
	27017: {Port: 27017, Service: "mongodb", Proto: "tcp"},
}

// New returns an Enricher with optional custom port mappings merged on top
// of the built-in registry.
func New(custom map[int]Entry) *Enricher {
	if custom == nil {
		custom = make(map[int]Entry)
	}
	return &Enricher{custom: custom}
}

// Lookup returns the Entry for the given port. If no mapping is found, a
// generic entry is returned with an empty service name.
func (e *Enricher) Lookup(port int) Entry {
	if entry, ok := e.custom[port]; ok {
		return entry
	}
	if entry, ok := wellKnown[port]; ok {
		return entry
	}
	return Entry{Port: port, Service: "unknown", Proto: "tcp"}
}

// Label returns a human-readable label such as "80/http" or "9999/unknown".
func (e *Enricher) Label(port int) string {
	en := e.Lookup(port)
	if en.Service == "unknown" {
		return fmt.Sprintf("%d", port)
	}
	return fmt.Sprintf("%d/%s", port, en.Service)
}
