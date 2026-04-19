// Package grouping provides port change grouping by service category.
package grouping

import (
	"sync"

	"github.com/user/portwatch/internal/state"
)

// Group represents a named collection of ports.
type Group struct {
	Name  string
	Ports []int
}

// Grouper maps port changes to named service groups.
type Grouper struct {
	mu     sync.RWMutex
	groups []Group
}

// New creates a Grouper with the provided groups.
func New(groups []Group) *Grouper {
	return &Grouper{groups: groups}
}

// Label returns the group name for a given port, or "unknown" if unmatched.
func (g *Grouper) Label(port int) string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, grp := range g.groups {
		for _, p := range grp.Ports {
			if p == port {
				return grp.Name
			}
		}
	}
	return "unknown"
}

// Annotate returns a map of port -> group label for each change.
func (g *Grouper) Annotate(changes []state.Change) map[int]string {
	out := make(map[int]string, len(changes))
	for _, c := range changes {
		out[c.Port] = g.Label(c.Port)
	}
	return out
}

// Add appends a new group at runtime.
func (g *Grouper) Add(grp Group) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.groups = append(g.groups, grp)
}
