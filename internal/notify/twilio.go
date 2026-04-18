package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/user/portwatch/internal/state"
)

// TwilioNotifier sends SMS alerts via the Twilio REST API.
type TwilioNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	baseURL    string
	client     *http.Client
}

// NewTwilioNotifier creates a Twilio SMS notifier.
func NewTwilioNotifier(accountSID, authToken, from, to string) *TwilioNotifier {
	return &TwilioNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		baseURL:    "https://api.twilio.com",
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches an SMS for the given port change.
func (t *TwilioNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	body := fmt.Sprintf("[portwatch] Port %d %s on %s", change.Port, action, change.Host)

	endpoint := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Messages.json", t.baseURL, t.accountSID)
	data := url.Values{}
	data.Set("From", t.from)
	data.Set("To", t.to)
	data.Set("Body", body)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("twilio: unexpected status %d", resp.StatusCode)
	}
	return nil
}
