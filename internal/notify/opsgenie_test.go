package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestOpsGenieNotifier_Send_OpenedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := &OpsGenieNotifier{apiKey: "test-key", apiURL: ts.URL, client: ts.Client()}
	err := n.Send(state.Change{Port: "8080", Type: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message"] != "Port 8080 opened" {
		t.Errorf("unexpected message: %v", received["message"])
	}
	if received["priority"] != "P2" {
		t.Errorf("expected P2 priority for opened port, got %v", received["priority"])
	}
}

func TestOpsGenieNotifier_Send_ClosedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := &OpsGenieNotifier{apiKey: "test-key", apiURL: ts.URL, client: ts.Client()}
	err := n.Send(state.Change{Port: "9090", Type: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["priority"] != "P3" {
		t.Errorf("expected P3 priority for closed port, got %v", received["priority"])
	}
}

func TestOpsGenieNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := &OpsGenieNotifier{apiKey: "bad-key", apiURL: ts.URL, client: ts.Client()}
	err := n.Send(state.Change{Port: "80", Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestOpsGenieNotifier_Send_BadURL(t *testing.T) {
	n := &OpsGenieNotifier{apiKey: "k", apiURL: "http://127.0.0.1:0", client: &http.Client{}}
	err := n.Send(state.Change{Port: "443", Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewOpsGenieNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewOpsGenieNotifier("key")
}
