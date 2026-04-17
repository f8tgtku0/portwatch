package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestMattermostNotifier_Send_OpenedPort(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL)
	err := n.Send(state.ChangeMessage{Change: state.Change{Port: 8080, Protocol: "tcp", Type: state.Opened}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text payload")
	}
}

func TestMattermostNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL)
	err := n.Send(state.ChangeMessage{Change: state.Change{Port: 22, Protocol: "tcp", Type: state.Closed}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMattermostNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewMattermostNotifier(ts.URL)
	err := n.Send(state.ChangeMessage{Change: state.Change{Port: 443, Protocol: "tcp", Type: state.Opened}})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMattermostNotifier_Send_BadURL(t *testing.T) {
	n := NewMattermostNotifier("http://127.0.0.1:0/bad")
	err := n.Send(state.ChangeMessage{Change: state.Change{Port: 80, Protocol: "tcp", Type: state.Opened}})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewMattermostNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewMattermostNotifier("http://example.com")
}
