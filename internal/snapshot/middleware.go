package snapshot

import (
	"log"

	"github.com/user/portwatch/internal/scanner"
)

// Recorder wraps a Snapshot and records each scan result automatically.
type Recorder struct {
	snap *Snapshot
}

// NewRecorder returns a Recorder backed by snap.
func NewRecorder(snap *Snapshot) *Recorder {
	return &Recorder{snap: snap}
}

// Apply records ports into the snapshot, logging any persistence error.
func (r *Recorder) Apply(ports []scanner.Port) []scanner.Port {
	if r.snap == nil {
		return ports
	}
	if err := r.snap.Record(ports); err != nil {
		log.Printf("snapshot: failed to record: %v", err)
	}
	return ports
}
