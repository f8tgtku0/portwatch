package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewGoogleChatNotifier returns a Notifier that posts messages to a Google Chat webhook.
func NewGoogleChatNotifier(webhookURL string) Notifier {
	return &googleChatNotifier{webhookURL: webhookURL}
}

type googleChatNotifier struct {
	webhookURL string
}

func (g *googleChatNotifier) Send(msg Message) error {
	action := "opened"
	if msg.Change.Type == state.Closed {
		action = "closed"
	}
	text := fmt.Sprintf("*Port %s*: port `%d` was %s", msg.Level, msg.Change.Port, action)

	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return fmt.Errorf("googlechat: marshal: %w", err)
	}

	resp, err := http.Post(g.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
