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

func jiraChange(opened bool) state.Change {
	return state.Change{Port: 8080, Opened: opened}
}

func TestJiraNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewJiraNotifier(ts.URL, "OPS", "Bug", "user@example.com", "token")
	if err := n.Send(jiraChange(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := got["fields"].(map[string]any)
	if fields["summary"] != "Port 8080 unexpectedly opened" {
		t.Errorf("unexpected summary: %v", fields["summary"])
	}
}

func TestJiraNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewJiraNotifier(ts.URL, "OPS", "Bug", "user@example.com", "token")
	if err := n.Send(jiraChange(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJiraNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewJiraNotifier(ts.URL, "OPS", "Bug", "user@example.com", "token")
	if err := n.Send(jiraChange(true)); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestJiraNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewJiraNotifier("http://127.0.0.1:0", "OPS", "Bug", "a", "b")
	if err := n.Send(jiraChange(true)); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewJiraNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewJiraNotifier("http://example.com", "P", "Bug", "e", "t")
}
