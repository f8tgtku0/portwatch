package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewMattermostNotifier sends alerts to a Mattermost incoming webhook.
func NewMattermostNotifier(webhookURL string) Notifier {
	return &mattermostNotifier{webhookURL: webhookURL}
}

type mattermostNotifier struct {
	webhookURL string
}

func (m *mattermostNotifier) Send(msg state.ChangeMessage) error {
	action := "closed"
	if msg.Change.Type == state.Opened {
		action = "opened"
	}
	text := fmt.Sprintf("**Port %s**: `%d` was **%s**", msg.Change.Protocol, msg.Change.Port, action)

	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return err
	}

	resp, err := http.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
