package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestSignalRNotifier_Send_OpenedPort(t *testing.T) {
	var received signalRPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSignalRNotifier(ts.URL, "portAlert")
	err := n.Send(history.Entry{Port: 8080, Proto: "tcp", Change: history.Opened, Timestamp: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Action != "opened" {
		t.Errorf("expected action 'opened', got %q", received.Action)
	}
	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Method != "portAlert" {
		t.Errorf("expected method 'portAlert', got %q", received.Method)
	}
}

func TestSignalRNotifier_Send_ClosedPort(t *testing.T) {
	var received signalRPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSignalRNotifier(ts.URL, "portAlert")
	err := n.Send(history.Entry{Port: 443, Proto: "tcp", Change: history.Closed, Timestamp: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Action != "closed" {
		t.Errorf("expected action 'closed', got %q", received.Action)
	}
}

func TestSignalRNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSignalRNotifier(ts.URL, "portAlert")
	err := n.Send(history.Entry{Port: 80, Proto: "tcp", Change: history.Opened, Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSignalRNotifier_Send_BadURL(t *testing.T) {
	n := NewSignalRNotifier("http://127.0.0.1:0/invalid", "portAlert")
	err := n.Send(history.Entry{Port: 80, Proto: "tcp", Change: history.Opened, Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewSignalRNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewSignalRNotifier("http://localhost", "method")
}
