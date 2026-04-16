package notify

import (
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

// minimalSMTPServer accepts one connection, performs a bare SMTP handshake,
// and records the raw data it receives.
func startFakeSMTP(t *testing.T) (addr string, received chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		ln.Close()
		if err != nil {
			ch <- ""
			return
		}
		defer conn.Close()
		tc := textproto.NewConn(conn)
		_ = tc // just drain
		buf, _ := io.ReadAll(conn)
		ch <- string(buf)
	}()
	return ln.Addr().String(), ch
}

func TestEmailNotifier_ActionLabel(t *testing.T) {
	if actionLabel(state.Change{Opened: true}) != "opened" {
		t.Error("expected opened")
	}
	if actionLabel(state.Change{Opened: false}) != "closed" {
		t.Error("expected closed")
	}
}

func TestEmailNotifier_Send_InvalidHost(t *testing.T) {
	cfg := EmailConfig{
		Host:     "127.0.0.1",
		Port:     1,
		Username: "u",
		Password: "p",
		From:     "a@example.com",
		To:       []string{"b@example.com"},
	}
	n := NewEmailNotifier(cfg)
	msg := Message{
		Change:    state.Change{Port: 8080, Proto: "tcp", Opened: true, Host: "localhost"},
		Timestamp: time.Now(),
	}
	err := n.Send(msg)
	if err == nil {
		t.Error("expected error connecting to invalid SMTP host")
	}
}

func TestNewEmailNotifier_ImplementsNotifier(t *testing.T) {
	cfg := EmailConfig{Host: "localhost", Port: 25, From: "x@x.com", To: []string{"y@y.com"}}
	var _ Notifier = NewEmailNotifier(cfg)
	_ = smtp.PlainAuth // ensure smtp import used
}
