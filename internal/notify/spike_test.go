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
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSpikeNotifierWithURL(ts.URL)
	err := n.Send(state.Change{Port: 8080, Type: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["status"] != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", got["status"])
	}
}

func TestSpikeNotifier_Send_ClosedPort(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewSpikeNotifierWithURL(ts.URL)
	err := n.Send(state.Change{Port: 9090, Type: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["status"] != "RESOLVED" {
		t.Errorf("expected RESOLVED, got %s", got["status"])
	}
}

func TestSpikeNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewSpikeNotifierWithURL(ts.URL)
	err := n.Send(state.Change{Port: 80, Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSpikeNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewSpikeNotifierWithURL("http://127.0.0.1:0")
	err := n.Send(state.Change{Port: 80, Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewSpikeNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewSpikeNotifier("test-key")
}
