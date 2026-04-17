package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/user/portwatch/internal/history"
)

// NewSMSNotifier creates a Twilio-based SMS notifier.
func NewSMSNotifier(accountSID, authToken, from, to string) *SMSNotifier {
	return &SMSNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// SMSNotifier sends port change alerts via Twilio SMS.
type SMSNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	client     *http.Client
}

func (s *SMSNotifier) Send(e history.Entry) error {
	action := "closed"
	if e.Change == "opened" {
		action = "opened"
	}
	body := fmt.Sprintf("[portwatch] Port %d %s on %s", e.Port, action, e.Time.Format(time.RFC822))

	data := url.Values{}
	data.Set("From", s.from)
	data.Set("To", s.to)
	data.Set("Body", body)

	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.accountSID)
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("sms: build request: %w", err)
	}
	req.SetBasicAuth(s.accountSID, s.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sms: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sms: unexpected status %d", resp.StatusCode)
	}
	return nil
}
