package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func TestPagerDutyNotifier_Send_OpenedPort(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(202)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key")
	n.url = ts.URL

	msg := Message{
		Text:      "port 8080 opened",
		Timestamp: time.Now(),
		Change:    state.Change{Type: state.ChangeTypeOpened, Port: 8080},
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "test-key" {
		t.Errorf("expected routing key 'test-key', got %q", received.RoutingKey)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("expected severity 'warning', got %q", received.Payload.Severity)
	}
	if received.Payload.Summary != msg.Text {
		t.Errorf("expected summary %q, got %q", msg.Text, received.Payload.Summary)
	}
}

func TestPagerDutyNotifier_Send_ClosedPort(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(202)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key2")
	n.url = ts.URL

	msg := Message{
		Text:      "port 9090 closed",
		Timestamp: time.Now(),
		Change:    state.Change{Type: state.ChangeTypeClosed, Port: 9090},
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Payload.Severity != "info" {
		t.Errorf("expected severity 'info', got %q", received.Payload.Severity)
	}
}

func TestPagerDutyNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key")
	n.url = ts.URL

	err := n.Send(Message{Timestamp: time.Now(), Change: state.Change{Type: state.ChangeTypeOpened}})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPagerDutyNotifier_Send_BadURL(t *testing.T) {
	n := NewPagerDutyNotifier("key")
	n.url = "http://127.0.0.1:0"

	err := n.Send(Message{Timestamp: time.Now(), Change: state.Change{Type: state.ChangeTypeOpened}})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewPagerDutyNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewPagerDutyNotifier("key")
}
