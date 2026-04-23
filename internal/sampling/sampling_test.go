package sampling_test

import (
	"testing"

	"github.com/user/portwatch/internal/sampling"
	"github.com/user/portwatch/internal/state"
)

func makeChanges(ports ...int) []state.Change {
	changes := make([]state.Change, len(ports))
	for i, p := range ports {
		changes[i] = state.Change{Port: p, Action: state.Opened}
	}
	return changes
}

func TestNew_ClampsRateAboveOne(t *testing.T) {
	s := sampling.New(5.0)
	if s.Rate() != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", s.Rate())
	}
}

func TestNew_ClampsRateBelowZero(t *testing.T) {
	s := sampling.New(-1.0)
	if s.Rate() <= 0 {
		t.Fatalf("expected positive rate, got %f", s.Rate())
	}
}

func TestNew_RateOne_AllowsAll(t *testing.T) {
	s := sampling.New(1.0)
	for i := 0; i < 50; i++ {
		if !s.Allow() {
			t.Fatal("rate=1.0 should always allow")
		}
	}
}

func TestSample_RateOne_ReturnsAll(t *testing.T) {
	s := sampling.New(1.0)
	changes := makeChanges(80, 443, 8080)
	got := s.Sample(changes)
	if len(got) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(got))
	}
}

func TestSample_EmptyInput_ReturnsNil(t *testing.T) {
	s := sampling.New(1.0)
	got := s.Sample(nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestSample_RateZeroPoint01_ReducesVolume(t *testing.T) {
	s := sampling.New(0.01)
	changes := makeChanges(make([]int, 1000)...)
	for i := range changes {
		changes[i].Port = i + 1
	}
	got := s.Sample(changes)
	// With rate 0.01 over 1000 samples we expect far fewer than 1000.
	if len(got) >= 500 {
		t.Fatalf("expected sampling to reduce volume significantly, got %d/1000", len(got))
	}
}

func TestMiddleware_Apply_NilSampler_ReturnsAll(t *testing.T) {
	m := sampling.NewMiddleware(nil)
	changes := makeChanges(22, 80)
	got := m.Apply(changes)
	if len(got) != len(changes) {
		t.Fatalf("expected %d, got %d", len(changes), len(got))
	}
}

func TestMiddleware_Apply_RateOne_ReturnsAll(t *testing.T) {
	m := sampling.NewMiddleware(sampling.New(1.0))
	changes := makeChanges(22, 80, 443)
	got := m.Apply(changes)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}
