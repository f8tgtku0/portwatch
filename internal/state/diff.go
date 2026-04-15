package state

import "github.com/user/portwatch/internal/scanner"

// ChangeType indicates whether a port was opened or closed.
type ChangeType string

const (
	Opened ChangeType = "opened"
	Closed ChangeType = "closed"
)

// Change represents a single port state transition.
type Change struct {
	Port int
	Kind ChangeType
}

// Compare returns the list of changes between two port snapshots.
// prev is the previously known set of open ports; curr is the current scan result.
func Compare(prev, curr []scanner.Port) []Change {
	prevSet := toSet(prev)
	currSet := toSet(curr)

	var changes []Change

	for port := range currSet {
		if !prevSet[port] {
			changes = append(changes, Change{Port: port, Kind: Opened})
		}
	}

	for port := range prevSet {
		if !currSet[port] {
			changes = append(changes, Change{Port: port, Kind: Closed})
		}
	}

	return changes
}

func toSet(ports []scanner.Port) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p.Number] = true
	}
	return s
}
