package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewBearyChat returns a Notifier that sends messages to a BearyChat incoming webhook.
func NewBearyChat(webhookURL string) Notifier {
	return &bearyChatNotifier{webhookURL: webhookURL}
}

type bearyChatNotifier struct {
	webhookURL string
}

func (n *bearyChatNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	payload := map[string]interface{}{
		"text": fmt.Sprintf("Port **%d** was %s on host `%s`.", change.Port, action, change.Host),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: marshal: %w", err)
	}
	resp, err := http.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// NewBearyChat alias for builder compatibility.
func NewBearyChat_Notifier(webhookURL string) Notifier {
	return NewBearyChat(webhookURL)
}
