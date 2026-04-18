package notify

import "net/http"

// NewClickUpNotifierWithURL creates a ClickUp notifier with a custom base URL for testing.
func NewClickUpNotifierWithURL(apiKey, listID, baseURL string, client *http.Client) Notifier {
	n := NewClickUpNotifier(apiKey, listID)
	c := n.(*clickUpNotifier)
	c.baseURL = baseURL
	if client != nil {
		c.client = client
	}
	return c
}
