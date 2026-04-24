// Package normalize provides port change normalization to ensure consistent
// representation of port numbers and actions before downstream processing.
package normalize

import "github.com/user/portwatch/internal/state"

// Normalizer canonicalizes port change slices by deduplicating entries
// with the same port+action pair and sorting by port number ascending.
type Normalizer struct{}

// New returns a new Normalizer.
func New() *Normalizer {
	return &Normalizer{}
}

// Apply deduplicates and sorts a slice of state.Change values.
// Duplicate entries (same Port and Action) are collapsed to a single entry.
// The returned slice is ordered by port number ascending, with opened changes
// before closed changes when ports are equal.
func (n *Normalizer) Apply(changes []state.Change) []state.Change {
	if len(changes) == 0 {
		return changes
	}

	seen := make(map[string]struct{}, len(changes))
	unique := make([]state.Change, 0, len(changes))

	for _, c := range changes {
		k := changeKey(c)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		unique = append(unique, c)
	}

	sortChanges(unique)
	return unique
}

// changeKey returns a string key unique to a port+action pair.
func changeKey(c state.Change) string {
	action := "opened"
	if c.Action == state.Closed {
		action = "closed"
	}
	return fmt.Sprintf("%d:%s", c.Port, action)
}

// sortChanges sorts changes in-place by port ascending, opened before closed.
func sortChanges(changes []state.Change) {
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Port != changes[j].Port {
			return changes[i].Port < changes[j].Port
		}
		// opened (0) before closed (1)
		return changes[i].Action < changes[j].Action
	})
}

import (
	"fmt"
	"sort"
)
