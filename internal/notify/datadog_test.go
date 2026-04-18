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

func datadogChange(t state.ChangeType) state.Change {
	return state.Change{Port: 9200, Host: "localhost", Type: t}
}

func TestDatadogNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	n := notify.NewDatadogNotifier("test-key", srv.URL)
	if err := n.Send(datadogChange(state.Opened)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["alert_type"] != "warning" {
		t.Errorf("expected alert_type warning, got %v", got["alert_type"])
	}
}

func TestDatadogNotifier_Send_ClosedPort(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	n := notify.NewDatadogNotifier("test-key", srv.URL)
	if err := n.Send(datadogChange(state.Closed)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["alert_type"] != "info" {
		t.Errorf("expected alert_type info, got %v", got["alert_type"])
	}
}

func TestDatadogNotifier_Send_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := notify.NewDatadogNotifier("bad-key", srv.URL)
	if err := n.Send(datadogChange(state.Opened)); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestDatadogNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewDatadogNotifier("key", "http://127.0.0.1:0")
	if err := n.Send(datadogChange(state.Opened)); err == nil {
		t.Error("expected error for bad URL")
	}
}

func TestNewDatadogNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewDatadogNotifier("key", "")
}
