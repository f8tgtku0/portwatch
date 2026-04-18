package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewFreshdeskNotifier creates a Notifier that opens a Freshdesk ticket on port changes.
func NewFreshdeskNotifier(domain, apiKey, email string) Notifier {
	base := fmt.Sprintf("https://%s.freshdesk.com/api/v2/tickets", domain)
	return &freshdeskNotifier{base: base, apiKey: apiKey, email: email, client: &http.Client{}}
}

type freshdeskNotifier struct {
	base   string
	apiKey string
	email  string
	client *http.Client
}

type freshdeskTicket struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Priority    int    `json:"priority"`
	Status      int    `json:"status"`
}

func (f *freshdeskNotifier) Send(c state.Change) error {
	action := "opened"
	priority := 2
	if c.Type == state.Closed {
		action = "closed"
		priority = 1
	}
	ticket := freshdeskTicket{
		Subject:     fmt.Sprintf("Port %d %s", c.Port, action),
		Description: fmt.Sprintf("portwatch detected port %d was %s on the monitored host.", c.Port, action),
		Email:       f.email,
		Priority:    priority,
		Status:      2,
	}
	body, err := json.Marshal(ticket)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, f.base, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(f.apiKey, "X")
	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshdesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
