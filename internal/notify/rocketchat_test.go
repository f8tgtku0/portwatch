package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestRocketChatNotifier_Send_OpenedPort(t *testing.T) {
	var received rocketChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	err := n.Send(state.Change{Type: state.Opened, Port: state.Port{Number: 8080, Proto: "tcp"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected non-empty text payload")
	}
}

func TestRocketChatNotifier_Send_ClosedPort(t *testing.T) {
	var received rocketChatPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	err := n.Send(state.Change{Type: state.Closed, Port: state.Port{Number: 22, Proto: "tcp"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRocketChatNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewRocketChatNotifier(ts.URL)
	err := n.Send(state.Change{Type: state.Opened, Port: state.Port{Number: 443, Proto: "tcp"}})
	if err == nil {
		t.Error("expected error for non-OK status")
	}
}

func TestRocketChatNotifier_Send_BadURL(t *testing.T) {
	n := NewRocketChatNotifier("http://127.0.0.1:0/bad")
	err := n.Send(state.Change{Type: state.Opened, Port: state.Port{Number: 80, Proto: "tcp"}})
	if err == nil {
		t.Error("expected error for bad URL")
	}
}

func TestNewRocketChatNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewRocketChatNotifier("http://example.com")
}
