package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestLarkNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Send(state.Change{Port: 8080, Proto: "tcp", Opened: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["msg_type"] != "text" {
		t.Errorf("expected msg_type=text, got %v", got["msg_type"])
	}
}

func TestLarkNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Send(state.Change{Port: 22, Proto: "tcp", Opened: false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLarkNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewLarkNotifier(ts.URL)
	if err := n.Send(state.Change{Port: 80, Proto: "tcp", Opened: true}); err == nil {
		t.Error("expected error for non-OK status")
	}
}

func TestLarkNotifier_Send_BadURL(t *testing.T) {
	n := NewLarkNotifier("http://127.0.0.1:0")
	if err := n.Send(state.Change{Port: 443, Proto: "tcp", Opened: true}); err == nil {
		t.Error("expected error for bad URL")
	}
}

func TestNewLarkNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewLarkNotifier("http://example.com")
}
