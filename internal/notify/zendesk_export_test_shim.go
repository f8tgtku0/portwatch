package notify

func NewZendeskNotifierWithURL(url, email, token string) Notifier {
	return &zendeskNotifier{url: url, email: email, token: token}
}
