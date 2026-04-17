package notify

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/history"
)

// NewNtfyNotifier sends notifications via ntfy.sh or a self-hosted ntfy server.
func NewNtfyNotifier(serverURL, topic string) Notifier {
	return &ntfyNotifier{serverURL: strings.TrimRight(serverURL, "/"), topic: topic}
}

type ntfyNotifier struct {
	serverURL string
	topic     string
}

func (n *ntfyNotifier) Send(e history.Entry) error {
	action := "opened"
	priority := "default"
	if e.Change == "closed" {
		action = "closed"
		priority = "high"
	}

	msg := fmt.Sprintf("Port %d was %s at %s", e.Port, action, e.Time.Format("15:04:05"))
	url := fmt.Sprintf("%s/%s", n.serverURL, n.topic)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(msg))
	if err != nil {
		return err
	}
	req.Header.Set("Title", fmt.Sprintf("Port %s: %d", action, e.Port))
	req.Header.Set("Priority", priority)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
