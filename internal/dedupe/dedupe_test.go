package dedupe_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/state"
)

func opened(port int) state.Change {
	return state.Change{Port: port, Type: state.Opened}
}

func closed(port int) state.Change {
	return state.Change{Port: port, Type: state.Closed}
}

func TestIsDuplicate_FirstCallNotDuplicate(t *testing.T) {
	d := dedupe.New(5 * time.Second)
	if d.IsDuplicate(opened(80)) {
		t.Fatal("first call should not be duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindowIsDuplicate(t *testing.T) {
	d := dedupe.New(5 * time.Second)
	d.IsDuplicate(opened(80))
	if !d.IsDuplicate(opened(80)) {
		t.Fatal("second call within window should be duplicate")
	}
}

func TestIsDuplicate_DifferentPortsAreIndependent(t *testing.T) {
	d := dedupe.New(5 * time.Second)
	d.IsDuplicate(opened(80))
	if d.IsDuplicate(opened(443)) {
		t.Fatal("different port should not be duplicate")
	}
}

func TestIsDuplicate_DifferentActionsAreIndependent(t *testing.T) {
	d := dedupe.New(5 * time.Second)
	d.IsDuplicate(opened(80))
	if d.IsDuplicate(closed(80)) {
		t.Fatal("different action on same port should not be duplicate")
	}
}

func TestFilter_RemovesDuplicates(t *testing.T) {
	d := dedupe.New(5 * time.Second)
	changes := []state.Change{opened(80), opened(443)}
	out := d.Filter(changes)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	// second pass — both are duplicates
	out2 := d.Filter(changes)
	if len(out2) != 0 {
		t.Fatalf("expected 0 after dedup, got %d", len(out2))
	}
}

func TestPrune_AllowsReentryAfterExpiry(t *testing.T) {
	d := dedupe.New(1 * time.Millisecond)
	d.IsDuplicate(opened(8080))
	time.Sleep(5 * time.Millisecond)
	d.Prune()
	if d.IsDuplicate(opened(8080)) {
		t.Fatal("after prune and expiry, should not be duplicate")
	}
}
