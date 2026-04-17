package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestPushoverNotifier_Send_OpenedPort(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":1}`))
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = ts.URL

	err := n.Send(state.Change{Port: 8080, Host: "localhost", Type: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody["token"] != "tok" {
		t.Errorf("expected token 'tok', got %q", gotBody["token"])
	}
	if gotBody["user"] != "usr" {
		t.Errorf("expected user 'usr', got %q", gotBody["user"])
	}
}

func TestPushoverNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":1}`))
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = ts.URL

	err := n.Send(state.Change{Port: 22, Host: "localhost", Type: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPushoverNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":0}`))
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = ts.URL

	err := n.Send(state.Change{Port: 80, Host: "localhost", Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for non-1 status")
	}
}

func TestPushoverNotifier_Send_BadURL(t *testing.T) {
	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = "http://127.0.0.1:0"

	err := n.Send(state.Change{Port: 80, Host: "localhost", Type: state.Opened})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewPushoverNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewPushoverNotifier("tok", "usr")
}
