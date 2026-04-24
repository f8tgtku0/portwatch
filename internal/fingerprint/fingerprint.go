// Package fingerprint generates stable identifiers for port scan results,
// enabling change detection across scan cycles without full state comparison.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// Fingerprinter computes a deterministic hash over a set of open ports.
type Fingerprinter struct{}

// New returns a new Fingerprinter.
func New() *Fingerprinter {
	return &Fingerprinter{}
}

// Compute returns a hex-encoded SHA-256 hash of the sorted port list.
// An empty port set returns the hash of an empty string.
func (f *Fingerprinter) Compute(ports []state.Port) string {
	sorted := make([]state.Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Proto != sorted[j].Proto {
			return sorted[i].Proto < sorted[j].Proto
		}
		return sorted[i].Number < sorted[j].Number
	})

	parts := make([]string, len(sorted))
	for i, p := range sorted {
		parts[i] = fmt.Sprintf("%s:%d", p.Proto, p.Number)
	}

	h := sha256.Sum256([]byte(strings.Join(parts, ",")))
	return hex.EncodeToString(h[:])
}

// Changed returns true when the fingerprint of current differs from prev.
func (f *Fingerprinter) Changed(prev, current []state.Port) bool {
	return f.Compute(prev) != f.Compute(current)
}
