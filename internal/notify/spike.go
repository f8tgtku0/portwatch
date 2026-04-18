package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewSpikeNotifier creates a Notifier that sends alerts to Spike.sh.
func NewSpikeNotifier(apiKey string) Notifier {
	return &spikeNotifier{
		apiKey:  apiKey,
		endpoint: "https://alert.spike.sh/" + apiKey,
	}
}

type spikeNotifier struct {
	apiKey   string
	endpoint string
}

type spikePayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func spikeStatus(c state.Change) string {
	if c.Type == state.Opened {
		return "CRITICAL"
	}
	return "RESOLVED"
}

func (s *spikeNotifier) Send(c state.Change) error {
	action := "opened"
	if c.Type == state.Closed {
		action = "closed"
	}
	payload := spikePayload{
		Title:   fmt.Sprintf("Port %d %s", c.Port, action),
		Message: fmt.Sprintf("Port %d was %s on the monitored host.", c.Port, action),
		Status:  spikeStatus(c),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(s.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("spike: unexpected status %d", resp.StatusCode)
	}
	return nil
}
