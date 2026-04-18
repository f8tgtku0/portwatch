package notify

// NewSignalRNotifierWithURL allows tests to inject a custom URL.
func NewSignalRNotifierWithURL(url, hub, method string) Notifier {
	n := NewSignalRNotifier(url, hub, method)
	return n
}
