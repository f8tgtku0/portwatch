package filter

import (
	"testing"
)

func TestNew_ValidSinglePort(t *testing.T) {
	f, err := New([]string{"8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Ignored(8080) {
		t.Error("expected 8080 to be ignored")
	}
	if f.Ignored(9090) {
		t.Error("expected 9090 not to be ignored")
	}
}

func TestNew_ValidRange(t *testing.T) {
	f, err := New([]string{"8000-8100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{8000, 8050, 8100} {
		if !f.Ignored(p) {
			t.Errorf("expected port %d to be ignored", p)
		}
	}
	if f.Ignored(7999) || f.Ignored(8101) {
		t.Error("ports outside range should not be ignored")
	}
}

func TestNew_InvalidSpec(t *testing.T) {
	_, err := New([]string{"notaport"})
	if err == nil {
		t.Error("expected error for invalid spec")
	}
}

func TestNew_InvalidRange_Reversed(t *testing.T) {
	_, err := New([]string{"9000-8000"})
	if err == nil {
		t.Error("expected error for reversed range")
	}
}

func TestNew_EmptySpecs(t *testing.T) {
	f, err := New([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Ignored(80) {
		t.Error("no rules, nothing should be ignored")
	}
}

func TestIgnored_MultipleRules(t *testing.T) {
	f, err := New([]string{"22", "8000-8080", "443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{22, 443, 8000, 8040, 8080} {
		if !f.Ignored(p) {
			t.Errorf("expected port %d to be ignored", p)
		}
	}
	if f.Ignored(80) {
		t.Error("port 80 should not be ignored")
	}
}
