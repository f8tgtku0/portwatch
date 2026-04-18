package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func linearChange(opened bool) state.Change {
	return state.Change{Port: 9090, Opened: opened}
}

func TestLinearNotifier_Send_OpenedPort(t *testing.T) {
	var body map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"issueCreate":{"success":true}}}`))
	}))
	defer ts.Close()

	n := notify.NewLinearNotifierWithURL("key", "team1", ts.URL)
	if err := n.Send(linearChange(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body["query"], "opened") {
		t.Errorf("expected 'opened' in query, got: %s", body["query"])
	}
}

func TestLinearNotifier_Send_ClosedPort(t *testing.T) {
	var body map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"issueCreate":{"success":true}}}`))
	}))
	defer ts.Close()

	n := notify.NewLinearNotifierWithURL("key", "team1", ts.URL)
	if err := n.Send(linearChange(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body["query"], "closed") {
		t.Errorf("expected 'closed' in query, got: %s", body["query"])
	}
}

func TestLinearNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewLinearNotifierWithURL("key", "team1", ts.URL)
	if err := n.Send(linearChange(true)); err == nil {
		t.Error("expected error on non-OK status")
	}
}

func TestLinearNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewLinearNotifierWithURL("key", "team1", "http://127.0.0.1:0")
	if err := n.Send(linearChange(true)); err == nil {
		t.Error("expected error on bad URL")
	}
}

func TestNewLinearNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewLinearNotifier("key", "team1")
}
