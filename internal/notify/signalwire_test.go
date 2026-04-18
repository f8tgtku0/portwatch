package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func signalWireChange(open bool) state.PortChange {
	return state.PortChange{Port: 8080, Host: "localhost", IsOpen: open}
}

func TestSignalWireNotifier_Send_OpenedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if body := r.FormValue("Body"); body == "" {
			t.Error("expected non-empty Body")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewSignalWireNotifier(ts.URL, "proj", "token", "+1000", "+2000")
	if err := n.Send(signalWireChange(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSignalWireNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewSignalWireNotifier(ts.URL, "proj", "token", "+1000", "+2000")
	if err := n.Send(signalWireChange(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSignalWireNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewSignalWireNotifier(ts.URL, "proj", "token", "+1000", "+2000")
	if err := n.Send(signalWireChange(true)); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSignalWireNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewSignalWireNotifier("http://127.0.0.1:0", "proj", "token", "+1000", "+2000")
	if err := n.Send(signalWireChange(true)); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestNewSignalWireNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewSignalWireNotifier("http://example.com", "p", "t", "+1", "+2")
}
