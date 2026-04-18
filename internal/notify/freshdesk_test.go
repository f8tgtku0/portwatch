package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func TestFreshdeskNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewFreshdeskNotifierWithURL(ts.URL, "key", "user@example.com")
	err := n.Send(state.Change{Port: 8080, Type: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got["subject"].(string), "8080") {
		t.Errorf("expected subject to contain port, got %q", got["subject"])
	}
	if !strings.Contains(got["subject"].(string), "opened") {
		t.Errorf("expected subject to mention opened, got %q", got["subject"])
	}
}

func TestFreshdeskNotifier_Send_ClosedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewFreshdeskNotifierWithURL(ts.URL, "key", "user@example.com")
	err := n.Send(state.Change{Port: 22, Type: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got["subject"].(string), "closed") {
		t.Errorf("expected subject to mention closed, got %q", got["subject"])
	}
}

func TestFreshdeskNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewFreshdeskNotifierWithURL(ts.URL, "key", "user@example.com")
	err := n.Send(state.Change{Port: 443, Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestFreshdeskNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewFreshdeskNotifierWithURL("http://127.0.0.1:1", "key", "user@example.com")
	err := n.Send(state.Change{Port: 80, Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewFreshdeskNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewFreshdeskNotifier("domain", "key", "user@example.com")
}
