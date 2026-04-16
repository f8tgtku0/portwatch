package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/history"
)

// SlackNotifier sends port change alerts to a Slack webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a Slack notifier targeting the given webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Send posts a Slack message for the given history entry.
func (s *SlackNotifier) Send(e history.Entry) error {
	action := "closed"
	if e.Change == "opened" {
		action = "opened :unlock:"
	}
	text := fmt.Sprintf("[portwatch] Port *%d* %s at %s", e.Port, action, e.Timestamp.Format(time.RFC3339))

	body, err := json.Marshal(slackPayload{Text: text})
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
