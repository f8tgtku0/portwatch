package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestGotifyNotifier_Send_OpenedPort(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "testtoken")
	err := n.Send(history.Entry{Port: 8080, Change: "opened", Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Priority != 5 {
		t.Errorf("expected priority 5, got %d", received.Priority)
	}
}

func TestGotifyNotifier_Send_ClosedPort(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "testtoken")
	err := n.Send(history.Entry{Port: 22, Change: "closed", Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Priority != 8 {
		t.Errorf("expected priority 8, got %d", received.Priority)
	}
}

func TestGotifyNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "badtoken")
	err := n.Send(history.Entry{Port: 80, Change: "opened", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestGotifyNotifier_Send_BadURL(t *testing.T) {
	n := NewGotifyNotifier("http://127.0.0.1:0", "token")
	err := n.Send(history.Entry{Port: 80, Change: "opened", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewGotifyNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewGotifyNotifier("http://localhost", "token")
}
