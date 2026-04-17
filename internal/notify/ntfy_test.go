package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestNtfyNotifier_Send_OpenedPort(t *testing.T) {
	var gotPriority string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	err := n.Send(history.Entry{Port: 8080, Change: "opened", Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPriority != "default" {
		t.Errorf("expected priority 'default', got %q", gotPriority)
	}
}

func TestNtfyNotifier_Send_ClosedPort(t *testing.T) {
	var gotPriority string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	err := n.Send(history.Entry{Port: 22, Change: "closed", Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPriority != "high" {
		t.Errorf("expected priority 'high', got %q", gotPriority)
	}
}

func TestNtfyNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	err := n.Send(history.Entry{Port: 80, Change: "opened", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestNtfyNotifier_Send_BadURL(t *testing.T) {
	n := NewNtfyNotifier("http://127.0.0.1:0", "portwatch")
	err := n.Send(history.Entry{Port: 80, Change: "opened", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewNtfyNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewNtfyNotifier("https://ntfy.sh", "portwatch")
}
