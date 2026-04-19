package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

// newTestServer creates a test HTTP server that captures the decoded discordPayload.
func newTestServer(t *testing.T, status int, received *discordPayload) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if received != nil {
			json.NewDecoder(r.Body).Decode(received)
		}
		w.WriteHeader(status)
	}))
}

func TestDiscordNotifier_Send_OpenedPort(t *testing.T) {
	var received discordPayload
	ts := newTestServer(t, http.StatusNoContent, &received)
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	err := n.Send(history.Entry{Port: 8080, Proto: "tcp", Change: "opened", Timestamp: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Embeds) == 0 {
		t.Fatal("expected embeds in payload")
	}
	if received.Embeds[0].Color != 0x2ecc71 {
		t.Errorf("expected green color for opened, got %x", received.Embeds[0].Color)
	}
}

func TestDiscordNotifier_Send_ClosedPort(t *testing.T) {
	var received discordPayload
	ts := newTestServer(t, http.StatusNoContent, &received)
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	err := n.Send(history.Entry{Port: 22, Proto: "tcp", Change: "closed", Timestamp: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Embeds[0].Color != 0xe74c3c {
		t.Errorf("expected red color for closed, got %x", received.Embeds[0].Color)
	}
}

func TestDiscordNotifier_Send_NonOKStatus(t *testing.T) {
	ts := newTestServer(t, http.StatusInternalServerError, nil)
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	err := n.Send(history.Entry{Port: 80, Proto: "tcp", Change: "opened", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestDiscordNotifier_Send_BadURL(t *testing.T) {
	n := NewDiscordNotifier("http://127.0.0.1:0/bad")
	err := n.Send(history.Entry{Port: 443, Proto: "tcp", Change: "opened", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewDiscordNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewDiscordNotifier("http://example.com")
}
