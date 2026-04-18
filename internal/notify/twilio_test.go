package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func twilioChange(opened bool) state.Change {
	return state.Change{Port: 22, Host: "localhost", Opened: opened, At: time.Now()}
}

func TestTwilioNotifier_Send_OpenedPort(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotBody = r.FormValue("Body")
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.(*notify.TwilioNotifier).SetBaseURL(ts.URL) // via export shim

	if err := n.Send(twilioChange(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody == "" {
		t.Error("expected non-empty SMS body")
	}
}

func TestTwilioNotifier_Send_ClosedPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	n := notify.NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.(*notify.TwilioNotifier).SetBaseURL(ts.URL)

	if err := n.Send(twilioChange(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTwilioNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := notify.NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.(*notify.TwilioNotifier).SetBaseURL(ts.URL)

	if err := n.Send(twilioChange(true)); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestTwilioNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewTwilioNotifier("ACtest", "token", "+10000000000", "+19999999999")
	n.(*notify.TwilioNotifier).SetBaseURL("http://127.0.0.1:0")
	if err := n.Send(twilioChange(true)); err == nil {
		t.Error("expected error for unreachable host")
	}
}

func TestNewTwilioNotifier_ImplementsNotifier(t *testing.T) {
	var _ notify.Notifier = notify.NewTwilioNotifier("a", "b", "c", "d")
}
