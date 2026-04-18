package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// NewSignalWireNotifier sends SMS alerts via the SignalWire REST API.
func NewSignalWireNotifier(spaceURL, projectID, apiToken, from, to string) Notifier {
	return &signalWireNotifier{
		spaceURL:  strings.TrimRight(spaceURL, "/"),
		projectID: projectID,
		apiToken:  apiToken,
		from:      from,
		to:        to,
		client:    &http.Client{},
	}
}

type signalWireNotifier struct {
	spaceURL  string
	projectID string
	apiToken  string
	from      string
	to        string
	client    *http.Client
}

func (s *signalWireNotifier) Send(change state.PortChange) error {
	action := "opened"
	if !change.IsOpen {
		action = "closed"
	}
	body := fmt.Sprintf("[portwatch] port %d %s on %s", change.Port, action, change.Host)

	endpoint := fmt.Sprintf("%s/api/laml/2010-04-01/Accounts/%s/Messages.json", s.spaceURL, s.projectID)

	form := url.Values{}
	form.Set("From", s.from)
	form.Set("To", s.to)
	form.Set("Body", body)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("signalwire: build request: %w", err)
	}
	req.SetBasicAuth(s.projectID, s.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalwire: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalwire: unexpected status %d", resp.StatusCode)
	}
	return nil
}
