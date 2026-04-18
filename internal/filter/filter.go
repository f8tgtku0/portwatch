// Package filter provides port filtering logic for portwatch.
// It allows users to ignore specific ports or ranges from monitoring.
package filter

import "fmt"

// Rule represents a single filter rule.
type Rule struct {
	Port  int
	Start int
	End   int
	range_ bool
}

// Filter holds a set of rules used to exclude ports from alerting.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a list of port specs.
// Specs can be single ports ("8080") or ranges ("8000-9000").
func New(specs []string) (*Filter, error) {
	f := &Filter{}
	for _, s := range specs {
		r, err := parseSpec(s)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid spec %q: %w", s, err)
		}
		f.rules = append(f.rules, r)
	}
	return f, nil
}

// Ignored returns true if the given port matches any filter rule.
func (f *Filter) Ignored(port int) bool {
	for _, r := range f.rules {
		if r.range_ {
			if port >= r.Start && port <= r.End {
				return true
			}
		} else {
			if port == r.Port {
				return true
			}
		}
	}
	return false
}

func parseSpec(s string) (Rule, error) {
	var start, end int
	if n, _ := fmt.Sscanf(s, "%d-%d", &start, &end); n == 2 {
		if start < 1 || end > 65535 || start > end {
			return Rule{}, fmt.Errorf("range out of bounds")
		}
		return Rule{Start: start, End: end, range_: true}, nil
	}
	var port int
	if n, _ := fmt.Sscanf(s, "%d", &port); n == 1 {
		if port < 1 || port > 65535 {
			return Rule{}, fmt.Errorf("port out of bounds")
		}
		return Rule{Port: port}, nil
	}
	return Rule{}, fmt.Errorf("unrecognised format")
}
