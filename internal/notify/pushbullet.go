package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

const defaultPushbulletURL = "https://api.pushbullet.com/v2/pushes"

type pushbulletNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

// NewPushbulletNotifier returns a Notifier that sends alerts via Pushbullet.
func NewPushbulletNotifier(apiKey string) Notifier {
	return &pushbulletNotifier{
		apiKey: apiKey,
		apiURL: defaultPushbulletURL,
		client: &http.Client{},
	}
}

func (p *pushbulletNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	payload := map[string]string{
		"type":  "note",
		"title": fmt.Sprintf("Port %s", action),
		"body":  fmt.Sprintf("Port %d was %s.", change.Port, action),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, p.apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", p.apiKey)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pushbullet: unexpected status %d", resp.StatusCode)
	}
	return nil
}
