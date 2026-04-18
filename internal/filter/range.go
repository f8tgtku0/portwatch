package filter

import "fmt"

// Range represents an inclusive port range.
type Range struct {
	Low  int
	High int
}

// NewRange creates a Range, validating that Low <= High and both are in [1,65535].
func NewRange(low, high int) (Range, error) {
	if low < 1 || high > 65535 {
		return Range{}, fmt.Errorf("port range %d-%d out of bounds [1,65535]", low, high)
	}
	if low > high {
		return Range{}, fmt.Errorf("port range %d-%d is reversed", low, high)
	}
	return Range{Low: low, High: high}, nil
}

// Contains reports whether port p falls within the range.
func (r Range) Contains(p int) bool {
	return p >= r.Low && p <= r.High
}

// String returns a human-readable representation of the range.
func (r Range) String() string {
	if r.Low == r.High {
		return fmt.Sprintf("%d", r.Low)
	}
	return fmt.Sprintf("%d-%d", r.Low, r.High)
}

// Overlaps reports whether r and other share at least one port.
func (r Range) Overlaps(other Range) bool {
	return r.Low <= other.High && other.Low <= r.High
}
