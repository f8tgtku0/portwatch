package notify

// NewSignalRNotifierWithURL allows tests to inject a custom URL.
// It returns a Notifier configured with the provided SignalR hub URL,
// hub name, and method name, bypassing any default URL resolution.
func NewSignalRNotifierWithURL(url, hub, method string) Notifier {
	n := NewSignalRNotifier(url, hub, method)
	return n
}

// NewSignalRNotifierWithClient allows tests to inject both a custom URL
// and a custom HTTP client for full control over outbound requests.
func NewSignalRNotifierWithClient(url, hub, method string, client HTTPClient) Notifier {
	n := NewSignalRNotifier(url, hub, method)
	if sr, ok := n.(*signalRNotifier); ok {
		sr.client = client
	}
	return n
}
