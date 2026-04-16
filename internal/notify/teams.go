package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/history"
)

// NewTeamsNotifier sends port change alerts to a Microsoft Teams incoming webhook.
func NewTeamsNotifier(webhookURL string) Notifier {
	return &teamsNotifier{webhookURL: webhookURL, client: &http.Client{Timeout: 5 * time.Second}}
}

type teamsNotifier struct {
	webhookURL string
	client     *http.Client
}

type teamsPayload struct {
	Type       string        `json:"@type"`
	Context    string        `json:"@context"`
	ThemeColor string        `json:"themeColor"`
	Summary    string        `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	ActivityText  string `json:"activityText"`
}

func (t *teamsNotifier) Send(msg history.Entry) error {
	color := "2ecc71"
	if msg.Change == "closed" {
		color = "e74c3c"
	}
	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: color,
		Summary:    fmt.Sprintf("Port %d/%s %s", msg.Port, msg.Proto, msg.Change),
		Sections: []teamsSection{
			{
				ActivityTitle: fmt.Sprintf("Port %s: %d/%s", msg.Change, msg.Port, msg.Proto),
				ActivityText:  fmt.Sprintf("Port **%d** (%s) was **%s** at %s", msg.Port, msg.Proto, msg.Change, msg.Timestamp.Format(time.RFC3339)),
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal: %w", err)
	}
	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
