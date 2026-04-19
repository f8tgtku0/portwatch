package notify

import "net/http"

func NewPushoverNotifierWithURL(apiKey, userKey, baseURL string, client *http.Client) Notifier {
	return &pushoverNotifier{
		apiKey:  apiKey,
		userKey: userKey,
		baseURL: baseURL,
		client:  client,
	}
}
