package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewZendeskNotifier sends port change alerts as Zendesk tickets.
func NewZendeskNotifier(subdomain, email, apiToken string) Notifier {
	baseURL := fmt.Sprintf("https://%s.zendesk.com/api/v2/tickets.json", subdomain)
	return &zendeskNotifier{url: baseURL, email: email, token: apiToken}
}

type zendeskNotifier struct {
	url   string
	email string
	token string
}

type zendeskTicket struct {
	Ticket zendeskTicketBody `json:"ticket"`
}

type zendeskTicketBody struct {
	Subject string          `json:"subject"`
	Comment zendeskComment  `json:"comment"`
	Priority string         `json:"priority"`
}

type zendeskComment struct {
	Body string `json:"body"`
}

func (z *zendeskNotifier) Send(c state.Change) error {
	action := "opened"
	priority := "normal"
	if c.Type == state.Closed {
		action = "closed"
		priority = "high"
	}

	payload := zendeskTicket{
		Ticket: zendeskTicketBody{
			Subject:  fmt.Sprintf("Port %s on %s", action, c.Host),
			Comment:  zendeskComment{Body: fmt.Sprintf("Port %d was %s on host %s.", c.Port, action, c.Host)},
			Priority: priority,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, z.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.email+"/token", z.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zendesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
