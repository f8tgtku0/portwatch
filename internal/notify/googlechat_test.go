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

func TestGoogleChatNotifier_Send_OpenedPort(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGoogleChatNotifier(ts.URL)
	err := n.Send(notify.Message{Level: "info", Change: state.Change{Port: 8080, Type: state.Opened}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody["text"] == "" {
		t.Error("expected non-empty text")
	}
}

func TestGoogleChatNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewGoogleChatNotifier(ts.URL)
	err := n.Send(notify.Message{Level: "warn", Change: state.Change{Port: 22, Type: state.Closed}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoogleChatNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewGoogleChatNotifier(ts.URL)
	err := n.Send(notify.Message{Level: "info", Change: state.Change{Port: 80, Type: state.Opened}})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestGoogleChatNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewGoogleChatNotifier("http://127.0.0.1:0")
	err := n.Send(notify.Message{Level: "info", Change: state.Change{Port: 443, Type: state.Opened}})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewGoogleChatNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewGoogleChatNotifier("http://example.com")
}
