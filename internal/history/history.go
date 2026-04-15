package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single port change event recorded in history.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Event     string    `json:"event"` // "opened" or "closed"
}

// History manages a rotating log of port change events.
type History struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
	filePath string
}

// New creates a new History instance backed by the given file path.
func New(filePath string, maxSize int) (*History, error) {
	h := &History{
		filePath: filePath,
		maxSize:  maxSize,
	}
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return h, nil
}

// Record appends a new entry to the history, trimming oldest if over maxSize.
func (h *History) Record(port int, event string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries = append(h.entries, Entry{
		Timestamp: time.Now().UTC(),
		Port:      port,
		Event:     event,
	})

	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}

	return h.save()
}

// Entries returns a copy of all recorded entries.
func (h *History) Entries() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]Entry, len(h.entries))
	copy(result, h.entries)
	return result
}

// Clear removes all entries and deletes the backing file.
func (h *History) Clear() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = nil
	return os.Remove(h.filePath)
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.filePath, data, 0644)
}

func (h *History) load() error {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}
