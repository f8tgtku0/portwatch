package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func TestSpikeNotifier_Send_OpenedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSpikeNotifier(ts.URL)
	if err := n.Send(state.Change{Port: 8080, Opened: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["status"] != "alert" {
		t.Errorf("expected status=alert, got %v", received["status"])
	}
}

func TestSpikeNotifier_Send_ClosedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSpikeNotifier(ts.URL)
	if err := n.Send(state.Change{Port: 22, Opened: false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["status"] != "resolved" {
		t.Errorf("expected status=resolved, got %v", received["status"])
	}
}

func TestSpikeNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewSpikeNotifier(ts.URL)
	if err := n.Send(state.Change{Port: 443, Opened: true}); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSpikeNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewSpikeNotifier("http://127.0.0.1:0/bad")
	if err := n.Send(state.Change{Port: 80, Opened: true}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewSpikeNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewSpikeNotifier("http://example.com")
}
