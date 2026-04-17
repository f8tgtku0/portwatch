package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewRocketChatNotifier returns a Notifier that posts messages to a Rocket.Chat webhook.
func NewRocketChatNotifier(webhookURL string) Notifier {
	return &rocketChatNotifier{webhookURL: webhookURL}
}

type rocketChatNotifier struct {
	webhookURL string
}

type rocketChatPayload struct {
	Text string `json:"text"`
}

func (r *rocketChatNotifier) Send(change state.Change) error {
	action := "closed"
	if change.Type == state.Opened {
		action = "opened"
	}

	msg := fmt.Sprintf(":bell: Port *%d* (%s) was *%s*", change.Port.Number, change.Port.Proto, action)

	payload, err := json.Marshal(rocketChatPayload{Text: msg})
	if err != nil {
		return fmt.Errorf("rocketchat: marshal payload: %w", err)
	}

	resp, err := http.Post(r.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("rocketchat: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
	}

	return nil
}
