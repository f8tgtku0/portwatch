package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestMatrixNotifier_Send_OpenedPort(t *testing.T) {
	var received matrixMessage
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Errorf("missing or wrong auth header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "!room:example.com", "tok")
	err := n.Send(Message{Title: "Port opened"}, state.Change{Port: 8080, Proto: "tcp", Action: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.MsgType != "m.text" {
		t.Errorf("expected msgtype m.text, got %s", received.MsgType)
	}
	if received.Body == "" {
		t.Error("expected non-empty body")
	}
}

func TestMatrixNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "!room:example.com", "tok")
	err := n.Send(Message{Title: "Port closed"}, state.Change{Port: 22, Proto: "tcp", Action: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMatrixNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "!room:example.com", "tok")
	err := n.Send(Message{Title: "Port opened"}, state.Change{Port: 9090, Proto: "tcp", Action: state.Opened})
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestMatrixNotifier_Send_BadURL(t *testing.T) {
	n := NewMatrixNotifier("http://127.0.0.1:0", "!room:example.com", "tok")
	err := n.Send(Message{Title: "Port opened"}, state.Change{Port: 80, Proto: "tcp", Action: state.Opened})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewMatrixNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewMatrixNotifier("http://localhost", "!r:h", "t")
}
