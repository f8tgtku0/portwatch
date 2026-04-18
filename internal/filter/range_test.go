package filter

import (
	"testing"
)

func TestNewRange_Valid(t *testing.T) {
	r, err := NewRange(80, 443)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Low != 80 || r.High != 443 {
		t.Errorf("expected 80-443, got %v", r)
	}
}

func TestNewRange_Reversed(t *testing.T) {
	_, err := NewRange(443, 80)
	if err == nil {
		t.Fatal("expected error for reversed range")
	}
}

func TestNewRange_OutOfBounds(t *testing.T) {
	_, err := NewRange(0, 100)
	if err == nil {
		t.Fatal("expected error for low=0")
	}
	_, err = NewRange(1, 70000)
	if err == nil {
		t.Fatal("expected error for high=70000")
	}
}

func TestRange_Contains(t *testing.T) {
	r, _ := NewRange(100, 200)
	if !r.Contains(150) {
		t.Error("expected 150 to be contained")
	}
	if r.Contains(99) {
		t.Error("expected 99 to not be contained")
	}
	if r.Contains(201) {
		t.Error("expected 201 to not be contained")
	}
}

func TestRange_String_Single(t *testing.T) {
	r, _ := NewRange(80, 80)
	if r.String() != "80" {
		t.Errorf("expected '80', got %q", r.String())
	}
}

func TestRange_String_Range(t *testing.T) {
	r, _ := NewRange(80, 443)
	if r.String() != "80-443" {
		t.Errorf("expected '80-443', got %q", r.String())
	}
}

func TestRange_Overlaps(t *testing.T) {
	a, _ := NewRange(100, 200)
	b, _ := NewRange(150, 250)
	c, _ := NewRange(201, 300)
	if !a.Overlaps(b) {
		t.Error("expected a and b to overlap")
	}
	if a.Overlaps(c) {
		t.Error("expected a and c to not overlap")
	}
}
