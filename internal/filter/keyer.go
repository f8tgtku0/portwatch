package filter

import "fmt"

// SpecFromPort returns a single-port spec string for the given port.
// Useful when programmatically constructing filter rules.
func SpecFromPort(port int) string {
	return fmt.Sprintf("%d", port)
}

// SpecFromRange returns a range spec string for the given inclusive bounds.
func SpecFromRange(start, end int) string {
	return fmt.Sprintf("%d-%d", start, end)
}
