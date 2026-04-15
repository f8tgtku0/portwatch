package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot represents a recorded port scan result at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Store persists and retrieves port snapshots to/from disk.
type Store struct {
	mu       sync.RWMutex
	filePath string
	last     *Snapshot
}

// New creates a new Store backed by the given file path.
func New(filePath string) *Store {
	return &Store{filePath: filePath}
}

// Save writes the current snapshot to disk and caches it in memory.
func (s *Store) Save(ports []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snap := &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return err
	}

	s.last = snap
	return nil
}

// Load reads the last snapshot from disk. Returns nil, nil if no snapshot exists yet.
func (s *Store) Load() (*Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.last != nil {
		return s.last, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}

	s.last = &snap
	return s.last, nil
}

// Clear removes the persisted snapshot file and resets in-memory cache.
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = nil
	err := os.Remove(s.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
