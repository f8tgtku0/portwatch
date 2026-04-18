package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

const datadogDefaultURL = "https://api.datadoghq.com/api/v1/events"

type datadogNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

type datadogEvent struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
	AlertType string `json:"alert_type"`
}

// NewDatadogNotifier creates a Notifier that sends events to Datadog.
func NewDatadogNotifier(apiKey, apiURL string) Notifier {
	if apiURL == "" {
		apiURL = datadogDefaultURL
	}
	return &datadogNotifier{apiKey: apiKey, apiURL: apiURL, client: &http.Client{}}
}

func (d *datadogNotifier) Send(change state.Change) error {
	action := "opened"
	alertType := "warning"
	if change.Type == state.Closed {
		action = "closed"
		alertType = "info"
	}

	event := datadogEvent{
		Title:     fmt.Sprintf("Port %s: %d", action, change.Port),
		Text:      fmt.Sprintf("Port %d was %s on %s.", change.Port, action, change.Host),
		Tags:      []string{fmt.Sprintf("port:%d", change.Port), fmt.Sprintf("action:%s", action)},
		AlertType: alertType,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("datadog: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("datadog: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status: %d", resp.StatusCode)
	}
	return nil
}
