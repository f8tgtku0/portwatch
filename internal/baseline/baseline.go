// Package baseline manages a trusted set of ports that are expected to be open.
// Ports in the baseline are excluded from alerting.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
)

// Baseline holds a set of trusted open ports.
type Baseline struct {
	mu    sync.RWMutex
	ports map[int]struct{}
	path  string
}

// New creates a Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{
		ports: make(map[int]struct{}),
		path:  path,
	}
}

// Load reads the baseline from disk. If the file does not exist, the baseline
// starts empty.
func (b *Baseline) Load() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	data, err := os.ReadFile(b.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var ports []int
	if err := json.Unmarshal(data, &ports); err != nil {
		return err
	}
	b.ports = make(map[int]struct{}, len(ports))
	for _, p := range ports {
		b.ports[p] = struct{}{}
	}
	return nil
}

// Save writes the current baseline to disk.
func (b *Baseline) Save() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ports := make([]int, 0, len(b.ports))
	for p := range b.ports {
		ports = append(ports, p)
	}
	data, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Add adds a port to the baseline.
func (b *Baseline) Add(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ports[port] = struct{}{}
}

// Remove removes a port from the baseline.
func (b *Baseline) Remove(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.ports, port)
}

// Contains reports whether port is in the baseline.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.ports[port]
	return ok
}

// Ports returns a snapshot of all baseline ports.
func (b *Baseline) Ports() []int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]int, 0, len(b.ports))
	for p := range b.ports {
		out = append(out, p)
	}
	return out
}
