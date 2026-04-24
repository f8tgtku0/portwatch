package normalize_test

import (
	"testing"

	"github.com/user/portwatch/internal/normalize"
	"github.com/user/portwatch/internal/state"
)

func opened(port int) state.Change {
	return state.Change{Port: port, Action: state.Opened}
}

func closed(port int) state.Change {
	return state.Change{Port: port, Action: state.Closed}
}

func TestApply_EmptyChanges_ReturnsEmpty(t *testing.T) {
	n := normalize.New()
	out := n.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d entries", len(out))
	}
}

func TestApply_DeduplicatesSamePortAndAction(t *testing.T) {
	n := normalize.New()
	input := []state.Change{opened(80), opened(80), opened(80)}
	out := n.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry after dedup, got %d", len(out))
	}
	if out[0].Port != 80 {
		t.Errorf("expected port 80, got %d", out[0].Port)
	}
}

func TestApply_KeepsDifferentActions(t *testing.T) {
	n := normalize.New()
	input := []state.Change{opened(443), closed(443)}
	out := n.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestApply_SortsByPortAscending(t *testing.T) {
	n := normalize.New()
	input := []state.Change{opened(9000), opened(80), opened(443)}
	out := n.Apply(input)
	ports := []int{out[0].Port, out[1].Port, out[2].Port}
	expected := []int{80, 443, 9000}
	for i, p := range expected {
		if ports[i] != p {
			t.Errorf("index %d: expected port %d, got %d", i, p, ports[i])
		}
	}
}

func TestApply_OpenedBeforeClosedForSamePort(t *testing.T) {
	n := normalize.New()
	input := []state.Change{closed(8080), opened(8080)}
	out := n.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Action != state.Opened {
		t.Errorf("expected opened first, got %v", out[0].Action)
	}
}

func TestApply_MixedDuplicatesAndUnique(t *testing.T) {
	n := normalize.New()
	input := []state.Change{
		opened(22), opened(22), closed(80), opened(443), closed(80),
	}
	out := n.Apply(input)
	// opened(22)x1, closed(80)x1, opened(443)x1 = 3
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}
