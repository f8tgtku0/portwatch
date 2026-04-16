package notify

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/user/portwatch/internal/state"
)

// EmailConfig holds SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type emailNotifier struct {
	cfg EmailConfig
}

// NewEmailNotifier creates a Notifier that sends alerts via SMTP.
func NewEmailNotifier(cfg EmailConfig) Notifier {
	return &emailNotifier{cfg: cfg}
}

func (e *emailNotifier) Send(msg Message) error {
	subject := fmt.Sprintf("[portwatch] Port %d %s", msg.Change.Port, actionLabel(msg.Change))
	body := fmt.Sprintf(
		"Time: %s\nHost: %s\nPort: %d\nProto: %s\nStatus: %s\n",
		msg.Timestamp.Format(time.RFC1123),
		msg.Change.Host,
		msg.Change.Port,
		msg.Change.Proto,
		actionLabel(msg.Change),
	)
	raw := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.cfg.From,
		strings.Join(e.cfg.To, ", "),
		subject,
		body,
	)
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, [](raw))
}

func actionLabel(c state.Change) string {
	if c.Opened {
		return "opened"
	}
	return "closed"
}
