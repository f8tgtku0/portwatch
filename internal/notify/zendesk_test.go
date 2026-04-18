package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func TestZendeskNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewZendeskNotifierWithURL(ts.URL, "user@example.com", "token123")
	err := n.Send(state.Change{Port: 8080, Host: "localhost", Type: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ticket := got["ticket"].(map[string]interface{})
	if ticket["priority"] != "normal" {
		t.Errorf("expected normal priority, got %v", ticket["priority"])
	}
}

func TestZendeskNotifier_Send_ClosedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewZendeskNotifierWithURL(ts.URL, "user@example.com", "token123")
	err := n.Send(state.Change{Port: 22, Host: "localhost", Type: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ticket := got["ticket"].(map[string]interface{})
	if ticket["priority"] != "high" {
		t.Errorf("expected high priority for closed port, got %v", ticket["priority"])
	}
}

func TestZendeskNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewZendeskNotifierWithURL(ts.URL, "u", "t")
	if err := n.Send(state.Change{Port: 80, Type: state.Opened}); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestZendeskNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewZendeskNotifierWithURL("http://127.0.0.1:0", "u", "t")
	if err := n.Send(state.Change{Port: 80, Type: state.Opened}); err == nil {
		t.Error("expected error for bad URL")
	}
}

func TestNewZendeskNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewZendeskNotifier("sub", "email", "token")
}
