package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestPushbulletNotifier_Send_OpenedPort(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := &pushbulletNotifier{apiKey: "key", apiURL: ts.URL, client: ts.Client()}
	err := n.Send(state.Change{Port: 8080, Opened: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody["title"] != "Port opened" {
		t.Errorf("expected title 'Port opened', got %q", gotBody["title"])
	}
}

func TestPushbulletNotifier_Send_ClosedPort(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := &pushbulletNotifier{apiKey: "key", apiURL: ts.URL, client: ts.Client()}
	err := n.Send(state.Change{Port: 22, Opened: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody["title"] != "Port closed" {
		t.Errorf("expected title 'Port closed', got %q", gotBody["title"])
	}
}

func TestPushbulletNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := &pushbulletNotifier{apiKey: "bad", apiURL: ts.URL, client: ts.Client()}
	err := n.Send(state.Change{Port: 443, Opened: true})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestPushbulletNotifier_Send_BadURL(t *testing.T) {
	n := &pushbulletNotifier{apiKey: "key", apiURL: "http://127.0.0.1:0", client: &http.Client{}}
	err := n.Send(state.Change{Port: 80, Opened: true})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewPushbulletNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewPushbulletNotifier("key")
}
