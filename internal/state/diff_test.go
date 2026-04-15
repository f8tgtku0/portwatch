package state

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func ports(nums ...int) []scanner.Port {
	var ps []scanner.Port
	for _, n := range nums {
		ps = append(ps, scanner.Port{Number: n})
	}
	return ps
}

func TestCompare_NoChanges(t *testing.T) {
	changes := Compare(ports(80, 443), ports(80, 443))
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestCompare_DetectsOpenedPort(t *testing.T) {
	changes := Compare(ports(80), ports(80, 8080))
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Port != 8080 || changes[0].Kind != Opened {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_DetectsClosedPort(t *testing.T) {
	changes := Compare(ports(80, 443), ports(80))
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Port != 443 || changes[0].Kind != Closed {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_EmptyPrev(t *testing.T) {
	changes := Compare(nil, ports(22, 80))
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Kind != Opened {
			t.Errorf("expected Opened, got %s for port %d", c.Kind, c.Port)
		}
	}
}

func TestCompare_EmptyCurr(t *testing.T) {
	changes := Compare(ports(22, 80), nil)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Kind != Closed {
			t.Errorf("expected Closed, got %s for port %d", c.Kind, c.Port)
		}
	}
}
