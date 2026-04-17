package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/history"
)

// NewGotifyNotifier sends notifications to a self-hosted Gotify server.
func NewGotifyNotifier(serverURL, token string) Notifier {
	return &gotifyNotifier{serverURL: serverURL, token: token}
}

type gotifyNotifier struct {
	serverURL string
	token     string
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

func (g *gotifyNotifier) Send(e history.Entry) error {
	action := "opened"
	priority := 5
	if e.Change == "closed" {
		action = "closed"
		priority = 8
	}

	p := gotifyPayload{
		Title:    fmt.Sprintf("Port %s: %d", action, e.Port),
		Message:  fmt.Sprintf("Port %d was %s on %s", e.Port, action, e.Time.Format("2006-01-02 15:04:05")),
		Priority: priority,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/message?token=%s", g.serverURL, g.token)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
