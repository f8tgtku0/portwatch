// Package labeler attaches human-readable service labels to port change
// messages based on IANA well-known assignments and user-defined overrides.
package labeler

import "github.com/user/portwatch/internal/state"

// wellKnown maps common port numbers to their canonical service names.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	9200: "elasticsearch",
	27017: "mongodb",
}

// Labeler resolves a service name for a given port number.
type Labeler struct {
	overrides map[int]string
}

// New returns a Labeler with optional user-defined overrides. Overrides take
// precedence over the built-in well-known table.
func New(overrides map[int]string) *Labeler {
	o := make(map[int]string, len(overrides))
	for k, v := range overrides {
		o[k] = v
	}
	return &Labeler{overrides: o}
}

// Lookup returns the service label for port, or an empty string if unknown.
func (l *Labeler) Lookup(port int) string {
	if l != nil {
		if name, ok := l.overrides[port]; ok {
			return name
		}
	}
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return ""
}

// Annotate adds a Label field to each change in the slice. Changes are
// returned unmodified if the Labeler is nil.
func (l *Labeler) Annotate(changes []state.Change) []state.Change {
	if l == nil {
		return changes
	}
	out := make([]state.Change, len(changes))
	for i, c := range changes {
		c.Label = l.Lookup(c.Port)
		out[i] = c
	}
	return out
}
