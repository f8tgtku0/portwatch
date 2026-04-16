package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookNotifier_Send_OpenedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, time.Second)
	err := n.Send(Message{Port: 8080, Proto: "tcp", Opened: true, Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["event"] != "opened" {
		t.Errorf("expected event=opened, got %v", received["event"])
	}
	if int(received["port"].(float64)) != 8080 {
		t.Errorf("expected port=8080, got %v", received["port"])
	}
}

func TestWebhookNotifier_Send_ClosedPort(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, time.Second)
	err := n.Send(Message{Port: 443, Proto: "tcp", Opened: false, Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["event"] != "closed" {
		t.Errorf("expected event=closed, got %v", received["event"])
	}
}

func TestWebhookNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, time.Second)
	err := n.Send(Message{Port: 80, Proto: "tcp", Opened: true, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestWebhookNotifier_Send_BadURL(t *testing.T) {
	n := NewWebhookNotifier("http://127.0.0.1:0", 200*time.Millisecond)
	err := n.Send(Message{Port: 22, Proto: "tcp", Opened: true, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
