package prestige_test

import (
	"testing"

	"github.com/user/portwatch/internal/prestige"
	"github.com/user/portwatch/internal/state"
)

func opened(port int) state.Change { return state.Change{Port: port, Action: state.Opened} }
func closed(port int) state.Change { return state.Change{Port: port, Action: state.Closed} }

func TestRecord_TracksOpens(t *testing.T) {
	tr := prestige.New()
	tr.Record([]state.Change{opened(80), opened(80), opened(443)})

	s, ok := tr.Get(80)
	if !ok {
		t.Fatal("expected score for port 80")
	}
	if s.Opens != 2 {
		t.Fatalf("expected 2 opens, got %d", s.Opens)
	}
}

func TestRecord_TracksCloses(t *testing.T) {
	tr := prestige.New()
	tr.Record([]state.Change{opened(22), closed(22)})

	s, _ := tr.Get(22)
	if s.Opens != 1 || s.Closes != 1 {
		t.Fatalf("expected 1 open and 1 close, got %+v", s)
	}
}

func TestRecord_EmptyChanges_NoOp(t *testing.T) {
	tr := prestige.New()
	tr.Record(nil)
	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker after nil record")
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	tr := prestige.New()
	_, ok := tr.Get(9999)
	if ok {
		t.Fatal("expected false for unknown port")
	}
}

func TestAll_ReturnsCopyOfAllScores(t *testing.T) {
	tr := prestige.New()
	tr.Record([]state.Change{opened(80), opened(443), opened(8080)})

	all := tr.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(all))
	}
}

func TestReset_ClearsScores(t *testing.T) {
	tr := prestige.New()
	tr.Record([]state.Change{opened(80)})
	tr.Reset()

	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker after reset")
	}
	if _, ok := tr.Get(80); ok {
		t.Fatal("expected port 80 to be cleared after reset")
	}
}

func TestRecord_AccumulatesAcrossMultipleCalls(t *testing.T) {
	tr := prestige.New()
	tr.Record([]state.Change{opened(80)})
	tr.Record([]state.Change{opened(80)})
	tr.Record([]state.Change{closed(80)})

	s, _ := tr.Get(80)
	if s.Opens != 2 {
		t.Fatalf("expected 2 opens, got %d", s.Opens)
	}
	if s.Closes != 1 {
		t.Fatalf("expected 1 close, got %d", s.Closes)
	}
}
