// Package snapshot provides periodic port-state snapshots for drift detection.
package snapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry is a single snapshot record.
type Entry struct {
	TakenAt time.Time      `json:"taken_at"`
	Ports   []scanner.Port `json:"ports"`
}

// Snapshot stores the most recent port snapshot and persists it to disk.
type Snapshot struct {
	mu   sync.RWMutex
	path string
	last *Entry
}

// New creates a Snapshot backed by the given file path.
func New(path string) *Snapshot {
	return &Snapshot{path: path}
}

// Record saves the current port list as the latest snapshot.
func (s *Snapshot) Record(ports []scanner.Port) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := &Entry{TakenAt: time.Now(), Ports: ports}
	s.last = e
	return s.persist(e)
}

// Last returns the most recently recorded snapshot, or nil.
func (s *Snapshot) Last() *Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.last
}

// Load reads the snapshot from disk into memory.
func (s *Snapshot) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	var e Entry
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return err
	}
	s.last = &e
	return nil
}

// Clear removes the snapshot file and resets in-memory state.
func (s *Snapshot) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = nil
	if err := os.Remove(s.path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *Snapshot) persist(e *Entry) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}
