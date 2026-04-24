package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/state"
)

func ports(nums ...int) []state.Port {
	ps := make([]state.Port, len(nums))
	for i, n := range nums {
		ps[i] = state.Port{Number: n, Proto: "tcp"}
	}
	return ps
}

func TestCompute_EmptyPorts(t *testing.T) {
	fp := fingerprint.New()
	h := fp.Compute(nil)
	if h == "" {
		t.Fatal("expected non-empty hash for nil ports")
	}
}

func TestCompute_StableAcrossCalls(t *testing.T) {
	fp := fingerprint.New()
	a := fp.Compute(ports(80, 443, 22))
	b := fp.Compute(ports(80, 443, 22))
	if a != b {
		t.Errorf("expected identical hashes, got %q and %q", a, b)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	fp := fingerprint.New()
	a := fp.Compute(ports(80, 443, 22))
	b := fp.Compute(ports(22, 80, 443))
	if a != b {
		t.Errorf("expected order-independent hashes, got %q and %q", a, b)
	}
}

func TestCompute_DifferentPortsDifferentHash(t *testing.T) {
	fp := fingerprint.New()
	a := fp.Compute(ports(80))
	b := fp.Compute(ports(443))
	if a == b {
		t.Error("expected different hashes for different port sets")
	}
}

func TestChanged_ReturnsTrueOnDiff(t *testing.T) {
	fp := fingerprint.New()
	if !fp.Changed(ports(80), ports(80, 443)) {
		t.Error("expected Changed to return true when ports differ")
	}
}

func TestChanged_ReturnsFalseOnSame(t *testing.T) {
	fp := fingerprint.New()
	if fp.Changed(ports(80, 443), ports(443, 80)) {
		t.Error("expected Changed to return false for same logical set")
	}
}

func TestChanged_BothEmpty(t *testing.T) {
	fp := fingerprint.New()
	if fp.Changed(nil, nil) {
		t.Error("expected Changed to return false for two empty sets")
	}
}
