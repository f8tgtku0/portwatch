package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestTelegramNotifier_Send_OpenedPort(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("token123", "chat456")
	n.baseURL = ts.URL

	err := n.Send(history.Entry{Port: 8080, Change: "opened", Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["chat_id"] != "chat456" {
		t.Errorf("expected chat_id chat456, got %s", received["chat_id"])
	}
}

func TestTelegramNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("token123", "chat456")
	n.baseURL = ts.URL

	err := n.Send(history.Entry{Port: 22, Change: "closed", Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTelegramNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("badtoken", "chat456")
	n.baseURL = ts.URL

	err := n.Send(history.Entry{Port: 80, Change: "opened", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestTelegramNotifier_Send_BadURL(t *testing.T) {
	n := NewTelegramNotifier("token", "chat")
	n.baseURL = "http://127.0.0.1:0"

	err := n.Send(history.Entry{Port: 443, Change: "closed", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewTelegramNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewTelegramNotifier("t", "c")
}
