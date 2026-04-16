package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestSlackNotifier_Send_OpenedPort(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	e := history.Entry{Port: 8080, Change: "opened", Timestamp: time.Now()}
	if err := n.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Fatal("expected non-empty slack message")
	}
}

func TestSlackNotifier_Send_ClosedPort(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	e := history.Entry{Port: 22, Change: "closed", Timestamp: time.Now()}
	if err := n.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Fatal("expected non-empty slack message")
	}
}

func TestSlackNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	e := history.Entry{Port: 443, Change: "opened", Timestamp: time.Now()}
	if err := n.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSlackNotifier_Send_BadURL(t *testing.T) {
	n := NewSlackNotifier("http://127.0.0.1:0/invalid")
	e := history.Entry{Port: 80, Change: "opened", Timestamp: time.Now()}
	if err := n.Send(e); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestNewSlackNotifier_ImplementsNotifier(t *testing.T) {
	var _ interface{ Send(history.Entry) error } = NewSlackNotifier("http://example.com")
}
