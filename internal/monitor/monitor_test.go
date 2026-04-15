package monitor

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestMonitor_DetectsNewPort(t *testing.T) {
	var mu sync.Mutex
	var alerts []string

	alertFn := func(msg string) {
		mu.Lock()
		defer mu.Unlock()
		alerts = append(alerts, msg)
	}

	port, stop := startTCPServer(t)
	defer stop()

	m := New("127.0.0.1", port, port, 50*time.Millisecond, alertFn)

	// Seed initial state with no open ports.
	m.lastState = &PortState{OpenPorts: map[int]bool{}}

	if err := m.scan(); err != nil {
		t.Fatalf("scan error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d: %v", len(alerts), alerts)
	}
	expected := fmt.Sprintf("[ALERT] new open port detected: 127.0.0.1:%d", port)
	if alerts[0] != expected {
		t.Errorf("unexpected alert message: %q", alerts[0])
	}
}

func TestMonitor_DetectsClosedPort(t *testing.T) {
	var mu sync.Mutex
	var alerts []string

	alertFn := func(msg string) {
		mu.Lock()
		defer mu.Unlock()
		alerts = append(alerts, msg)
	}

	port, stop := startTCPServer(t)
	stop() // close immediately so port appears closed on next scan

	m := New("127.0.0.1", port, port, 50*time.Millisecond, alertFn)
	m.lastState = &PortState{OpenPorts: map[int]bool{port: true}}

	if err := m.scan(); err != nil {
		t.Fatalf("scan error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d: %v", len(alerts), alerts)
	}
	expected := fmt.Sprintf("[INFO] port closed: 127.0.0.1:%d", port)
	if alerts[0] != expected {
		t.Errorf("unexpected alert message: %q", alerts[0])
	}
}
