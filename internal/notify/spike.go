package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewSpikeNotifier creates a Notifier that sends alerts to Spike.sh via webhook.
func NewSpikeNotifier(webhookURL string) Notifier {
	return &spikeNotifier{webhookURL: webhookURL}
}

type spikeNotifier struct {
	webhookURL string
}

func (s *spikeNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}

	payload := map[string]interface{}{
		"title":   fmt.Sprintf("Port %d %s", change.Port, action),
		"message": fmt.Sprintf("Port %d was %s on the monitored host.", change.Port, action),
		"status":  spikeStatus(change.Opened),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("spike: marshal payload: %w", err)
	}

	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("spike: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("spike: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func spikeStatus(opened bool) string {
	if opened {
		return "alert"
	}
	return "resolved"
}
