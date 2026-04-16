package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestTeamsNotifier_Send_OpenedPort(t *testing.T) {
	var received teamsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	err := n.Send(history.Entry{Port: 8080, Proto: "tcp", Change: "opened", Timestamp: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.ThemeColor != "2ecc71" {
		t.Errorf("expected green theme for opened, got %s", received.ThemeColor)
	}
	if received.Type != "MessageCard" {
		t.Errorf("expected MessageCard type, got %s", received.Type)
	}
}

func TestTeamsNotifier_Send_ClosedPort(t *testing.T) {
	var received teamsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	err := n.Send(history.Entry{Port: 22, Proto: "tcp", Change: "closed", Timestamp: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.ThemeColor != "e74c3c" {
		t.Errorf("expected red theme for closed, got %s", received.ThemeColor)
	}
}

func TestTeamsNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	err := n.Send(history.Entry{Port: 80, Proto: "tcp", Change: "opened", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestTeamsNotifier_Send_BadURL(t *testing.T) {
	n := NewTeamsNotifier("http://127.0.0.1:0/bad")
	err := n.Send(history.Entry{Port: 443, Proto: "tcp", Change: "opened", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewTeamsNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewTeamsNotifier("http://example.com")
}
