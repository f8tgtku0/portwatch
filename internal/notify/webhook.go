package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends change notifications as JSON POST requests.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Event   string `json:"event"`
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
func NewWebhookNotifier(url string, timeout time.Duration) *WebhookNotifier {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &WebhookNotifier{
		url: url,
		client: &http.Client{Timeout: timeout},
	}
}

// Send posts the message as a JSON payload to the configured webhook URL.
func (w *WebhookNotifier) Send(msg Message) error {
	event := "opened"
	if !msg.Opened {
		event = "closed"
	}
	payload := webhookPayload{
		Event:   event,
		Port:    msg.Port,
		Proto:   msg.Proto,
		Message: fmt.Sprintf("port %d/%s %s", msg.Port, msg.Proto, event),
		Time:    msg.Time.UTC().Format(time.RFC3339),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
