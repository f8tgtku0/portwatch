package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func clickupChange(opened bool) state.Change {
	return state.Change{Port: 9200, Opened: opened}
}

func TestClickUpNotifier_Send_OpenedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := &ClickUpNotifier{token: "tok", listID: "123", baseURL: ts.URL}
	if err := n.Send(clickupChange(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["name"] != "Port 9200 opened" {
		t.Errorf("unexpected name: %v", got["name"])
	}
}

func TestClickUpNotifier_Send_ClosedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := &ClickUpNotifier{token: "tok", listID: "123", baseURL: ts.URL}
	if err := n.Send(clickupChange(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["name"] != "Port 9200 closed" {
		t.Errorf("unexpected name: %v", got["name"])
	}
}

func TestClickUpNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := &ClickUpNotifier{token: "tok", listID: "123", baseURL: ts.URL}
	if err := n.Send(clickupChange(true)); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestClickUpNotifier_Send_BadURL(t *testing.T) {
	n := &ClickUpNotifier{token: "tok", listID: "123", baseURL: "http://127.0.0.1:0"}
	if err := n.Send(clickupChange(true)); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewClickUpNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewClickUpNotifier("tok", "list1")
}
