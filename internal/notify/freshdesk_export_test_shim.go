package notify

import "net/http"

// NewFreshdeskNotifierWithURL creates a Freshdesk notifier with a custom base URL for testing.
func NewFreshdeskNotifierWithURL(baseURL, apiKey, email string) Notifier {
	return &freshdeskNotifier{base: baseURL, apiKey: apiKey, email: email, client: &http.Client{}}
}
