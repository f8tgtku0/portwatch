package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortState represents the last known state of scanned ports.
type PortState struct {
	OpenPorts map[int]bool
}

// Monitor continuously scans a port range and alerts on changes.
type Monitor struct {
	host     string
	startPort int
	endPort   int
	interval  time.Duration
	lastState *PortState
	scanner   *scanner.Scanner
	alertFn   func(msg string)
}

// New creates a new Monitor instance.
func New(host string, startPort, endPort int, interval time.Duration, alertFn func(string)) *Monitor {
	return &Monitor{
		host:      host,
		startPort: startPort,
		endPort:   endPort,
		interval:  interval,
		scanner:   scanner.New(host, 500*time.Millisecond),
		alertFn:   alertFn,
	}
}

// Start begins the monitoring loop. It blocks until the done channel is closed.
func (m *Monitor) Start(done <-chan struct{}) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	if err := m.scan(); err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}
	log.Printf("[portwatch] initial scan complete: %d open ports found", len(m.lastState.OpenPorts))

	for {
		select {
		case <-done:
			log.Println("[portwatch] monitor stopped")
			return nil
		case <-ticker.C:
			if err := m.scan(); err != nil {
				log.Printf("[portwatch] scan error: %v", err)
			}
		}
	}
}

func (m *Monitor) scan() error {
	ports, err := m.scanner.Scan(m.startPort, m.endPort)
	if err != nil {
		return err
	}

	current := &PortState{OpenPorts: make(map[int]bool, len(ports))}
	for _, p := range ports {
		current.OpenPorts[p] = true
	}

	if m.lastState != nil {
		m.diff(m.lastState, current)
	}

	m.lastState = current
	return nil
}

func (m *Monitor) diff(prev, curr *PortState) {
	for port := range curr.OpenPorts {
		if !prev.OpenPorts[port] {
			msg := fmt.Sprintf("[ALERT] new open port detected: %s:%d", m.host, port)
			m.alertFn(msg)
		}
	}
	for port := range prev.OpenPorts {
		if !curr.OpenPorts[port] {
			msg := fmt.Sprintf("[INFO] port closed: %s:%d", m.host, port)
			m.alertFn(msg)
		}
	}
}
