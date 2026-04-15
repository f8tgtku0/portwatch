package state

// Diff holds the result of comparing two port lists.
type Diff struct {
	Opened []int
	Closed []int
}

// HasChanges returns true if any ports were opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare returns the difference between a previous and current set of open ports.
// Opened contains ports present in current but not previous.
// Closed contains ports present in previous but not current.
func Compare(previous, current []int) Diff {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var opened, closed []int

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}

	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}

	return Diff{Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
