package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/state"
)

// TestFilter_WithDiff verifies that filter correctly removes ignored ports
// from a set of state changes before they would be acted upon.
func TestFilter_WithDiff(t *testing.T) {
	ports := []int{22, 80, 443, 8080, 9000}

	f, err := filter.New([]string{"22", "443", "8000-8100"})
	if err != nil {
		t.Fatalf("failed to build filter: %v", err)
	}

	var visible []int
	for _, p := range ports {
		if !f.Ignored(p) {
			visible = append(visible, p)
		}
	}

	expected := []int{80, 9000}
	if len(visible) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, visible)
	}
	for i, p := range expected {
		if visible[i] != p {
			t.Errorf("position %d: expected %d, got %d", i, p, visible[i])
		}
	}
	_ = state.New // ensure import is used in a real build context
}
