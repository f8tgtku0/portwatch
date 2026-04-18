package notify

// SetBaseURL allows tests to override the Twilio API base URL.
func (t *TwilioNotifier) SetBaseURL(u string) {
	t.baseURL = u
}
