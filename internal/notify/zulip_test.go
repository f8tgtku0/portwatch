package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestZulipNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewZulipNotifier(ts.URL, "bot@example.com", "key", "alerts", "ports")
	err := n.Send(state.Change{Port: 8080, Proto: "tcp", Opened: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["content"] == "" {
		t.Error("expected non-empty content")
	}
}

func TestZulipNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewZulipNotifier(ts.URL, "bot@example.com", "key", "alerts", "ports")
	if err := n.Send(state.Change{Port: 22, Proto: "tcp", Opened: false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestZulipNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewZulipNotifier(ts.URL, "bot@example.com", "key", "alerts", "ports")
	if err := n.Send(state.Change{Port: 80, Proto: "tcp", Opened: true}); err == nil {
		t.Error("expected error for non-OK status")
	}
}

func TestZulipNotifier_Send_BadURL(t *testing.T) {
	n := NewZulipNotifier("http://127.0.0.1:0", "bot@example.com", "key", "alerts", "ports")
	if err := n.Send(state.Change{Port: 443, Proto: "tcp", Opened: true}); err == nil {
		t.Error("expected error for bad URL")
	}
}

func TestNewZulipNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewZulipNotifier("http://example.com", "bot@example.com", "key", "alerts", "ports")
}
