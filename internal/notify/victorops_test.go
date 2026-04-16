package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestVictorOpsNotifier_Send_OpenedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL)
	err := n.Send(alert.Message{Port: 8080, Change: "opened", Level: alert.LevelWarning, Text: "port 8080 opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message_type"] != "WARNING" {
		t.Errorf("expected WARNING, got %v", received["message_type"])
	}
	if received["entity_id"] != "portwatch-port-8080" {
		t.Errorf("unexpected entity_id: %v", received["entity_id"])
	}
}

func TestVictorOpsNotifier_Send_ClosedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL)
	err := n.Send(alert.Message{Port: 22, Change: "closed", Level: alert.LevelCritical, Text: "port 22 closed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message_type"] != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %v", received["message_type"])
	}
}

func TestVictorOpsNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewVictorOpsNotifier(ts.URL)
	err := n.Send(alert.Message{Port: 80, Change: "opened", Level: alert.LevelInfo, Text: "port 80 opened"})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestVictorOpsNotifier_Send_BadURL(t *testing.T) {
	n := NewVictorOpsNotifier("http://127.0.0.1:0/invalid")
	err := n.Send(alert.Message{Port: 443, Change: "opened", Level: alert.LevelInfo, Text: "port 443 opened"})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewVictorOpsNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewVictorOpsNotifier("http://example.com")
}
