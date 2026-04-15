package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// Snapshot holds the result of a full port scan.
type Snapshot struct {
	Timestamp time.Time
	Ports     []PortState
}

// Scanner scans a range of TCP ports on a given host.
type Scanner struct {
	Host    string
	MinPort int
	MaxPort int
	Timeout time.Duration
}

// New creates a Scanner with sensible defaults.
func New(host string, minPort, maxPort int) *Scanner {
	return &Scanner{
		Host:    host,
		MinPort: minPort,
		MaxPort: maxPort,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan performs a TCP connect scan across the configured port range
// and returns a Snapshot of open ports.
func (s *Scanner) Scan() (*Snapshot, error) {
	if s.MinPort < 1 || s.MaxPort > 65535 || s.MinPort > s.MaxPort {
		return nil, fmt.Errorf("invalid port range %d-%d", s.MinPort, s.MaxPort)
	}

	var open []PortState

	for port := s.MinPort; port <= s.MaxPort; port++ {
		addr := net.JoinHostPort(s.Host, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err != nil {
			if !isRefused(err) {
				// non-refused errors (e.g. network unreachable) are soft-skipped
				continue
			}
			continue
		}
		conn.Close()
		open = append(open, PortState{Port: port, Protocol: "tcp", Open: true})
	}

	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     open,
	}, nil
}

// isRefused returns true when the error is a connection-refused response,
// which means the port is closed but the host is reachable.
func isRefused(err error) bool {
	return strings.Contains(err.Error(), "refused")
}
