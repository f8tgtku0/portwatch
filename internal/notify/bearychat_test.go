package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func bearyChatChange(opened bool) state.Change {
	return state.Change{Port: 8080, Host: "localhost", Opened: opened}
}

func TestBearyChat_Send_OpenedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewBearyChat(ts.URL)
	if err := n.Send(bearyChatChange(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text, _ := got["text"].(string)
	if text == "" {
		t.Error("expected non-empty text field")
	}
}

func TestBearyChat_Send_ClosedPort(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewBearyChat(ts.URL)
	if err := n.Send(bearyChatChange(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text, _ := got["text"].(string)
	if text == "" {
		t.Error("expected non-empty text field")
	}
}

func TestBearyChat_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewBearyChat(ts.URL)
	if err := n.Send(bearyChatChange(true)); err == nil {
		t.Error("expected error for non-OK status")
	}
}

func TestBearyChat_Send_BadURL(t *testing.T) {
	n := notify.NewBearyChat("http://127.0.0.1:0")
	if err := n.Send(bearyChatChange(true)); err == nil {
		t.Error("expected error for bad URL")
	}
}

func TestNewBearyChat_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewBearyChat("http://example.com")
}
