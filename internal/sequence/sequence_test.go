package sequence_test

import (
	"sync"
	"testing"

	"portwatch/internal/sequence"
	"portwatch/internal/state"
)

func opened(port int) state.Change {
	return state.Change{Port: port, Action: state.Opened}
}

func TestStamp_AssignsSequenceNumber(t *testing.T) {
	s := sequence.New()
	changes := []state.Change{opened(80)}

	stamped := s.Stamp(changes)

	if len(stamped) != 1 {
		t.Fatalf("expected 1 stamped change, got %d", len(stamped))
	}
	if stamped[0].Seq != 1 {
		t.Errorf("expected seq=1, got %d", stamped[0].Seq)
	}
}

func TestStamp_IncrementsPerBatch(t *testing.T) {
	s := sequence.New()

	first := s.Stamp([]state.Change{opened(80)})
	second := s.Stamp([]state.Change{opened(443)})

	if first[0].Seq != 1 {
		t.Errorf("expected first seq=1, got %d", first[0].Seq)
	}
	if second[0].Seq != 2 {
		t.Errorf("expected second seq=2, got %d", second[0].Seq)
	}
}

func TestStamp_SameBatchSharesSeq(t *testing.T) {
	s := sequence.New()
	changes := []state.Change{opened(80), opened(443), opened(8080)}

	stamped := s.Stamp(changes)

	for i, st := range stamped {
		if st.Seq != 1 {
			t.Errorf("change[%d]: expected seq=1, got %d", i, st.Seq)
		}
	}
}

func TestStamp_EmptyReturnsNil(t *testing.T) {
	s := sequence.New()
	result := s.Stamp(nil)
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

func TestCurrent_ReflectsLatest(t *testing.T) {
	s := sequence.New()
	if s.Current() != 0 {
		t.Errorf("expected initial current=0, got %d", s.Current())
	}
	s.Stamp([]state.Change{opened(80)})
	if s.Current() != 1 {
		t.Errorf("expected current=1 after one batch, got %d", s.Current())
	}
}

func TestReset_SetsCounterToZero(t *testing.T) {
	s := sequence.New()
	s.Stamp([]state.Change{opened(80)})
	s.Reset()
	if s.Current() != 0 {
		t.Errorf("expected current=0 after reset, got %d", s.Current())
	}
}

func TestStamp_ConcurrentSafe(t *testing.T) {
	s := sequence.New()
	var wg sync.WaitGroup
	const goroutines = 50
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s.Stamp([]state.Change{opened(80)})
		}()
	}
	wg.Wait()
	if s.Current() != goroutines {
		t.Errorf("expected current=%d, got %d", goroutines, s.Current())
	}
}
