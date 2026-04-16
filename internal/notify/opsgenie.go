package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/state"
)

const opsgenieAPIURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

// NewOpsGenieNotifier creates a new OpsGenieNotifier.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey: apiKey,
		apiURL: opsgenieAPIURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (o *OpsGenieNotifier) Send(change state.Change) error {
	action := "opened"
	priority := "P2"
	if change.Type == state.Closed {
		action = "closed"
		priority = "P3"
	}

	payload := map[string]interface{}{
		"message":  fmt.Sprintf("Port %s %s", change.Port, action),
		"alias":    fmt.Sprintf("portwatch-port-%s", change.Port),
		"priority": priority,
		"tags":     []string{"portwatch", action},
		"details": map[string]string{
			"port":   change.Port,
			"action": action,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
