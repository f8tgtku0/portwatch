package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/state"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyNotifier sends alerts to PagerDuty via the Events API v2.
type PagerDutyNotifier struct {
	integrationKey string
	client         *http.Client
	url            string
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewPagerDutyNotifier creates a PagerDutyNotifier with the given integration key.
func NewPagerDutyNotifier(integrationKey string) *PagerDutyNotifier {
	return &PagerDutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
		url:            pagerDutyEventURL,
	}
}

// Send delivers a port change notification to PagerDuty.
func (p *PagerDutyNotifier) Send(msg Message) error {
	action := "trigger"
	severity := "warning"
	if msg.Change.Type == state.ChangeTypeClosed {
		severity = "info"
	}

	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: action,
		Payload: pdDetails{
			Summary:   msg.Text,
			Source:    "portwatch",
			Severity:  severity,
			Timestamp: msg.Timestamp.UTC().Format(time.RFC3339),
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal: %w", err)
	}

	resp, err := p.client.Post(p.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
