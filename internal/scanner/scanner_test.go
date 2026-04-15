package scanner

import (
	"net"
	"testing"
	"time"
)

// startTCPServer opens a listener on a random OS-assigned port and returns
// the port number together with a cleanup function.
func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, cleanup := startTCPServer(t)
	defer cleanup()

	s := New("127.0.0.1", port, port)
	s.Timeout = 200 * time.Millisecond

	snap, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(snap.Ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(snap.Ports))
	}
	if snap.Ports[0].Port != port {
		t.Errorf("expected port %d, got %d", port, snap.Ports[0].Port)
	}
	if snap.Ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", snap.Ports[0].Protocol)
	}
	if snap.Timestamp.IsZero() {
		t.Error("snapshot timestamp should not be zero")
	}
}

func TestScan_NoOpenPorts(t *testing.T) {
	// Use a port range that is very unlikely to have listeners in CI.
	s := New("127.0.0.1", 60000, 60005)
	s.Timeout = 100 * time.Millisecond

	snap, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned unexpected error: %v", err)
	}
	// We can't guarantee these ports are closed, but we can verify the
	// snapshot is well-formed.
	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := New("127.0.0.1", 1000, 500) // min > max
	_, err := s.Scan()
	if err == nil {
		t.Error("expected error for invalid port range, got nil")
	}
}
