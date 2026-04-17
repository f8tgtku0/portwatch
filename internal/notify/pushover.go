package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/state"
)

const pushoverAPI = "https://api.pushover.net/1/messages.json"

type PushoverNotifier struct {
	token   string
	user    string
	apiURL  string
	client  *http.Client
}

func NewPushoverNotifier(token, user string) *PushoverNotifier {
	return &PushoverNotifier{
		token:  token,
		user:   user,
		apiURL: pushoverAPI,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *PushoverNotifier) Send(change state.Change) error {
	action := "opened"
	if change.Type == state.Closed {
		action = "closed"
	}
	message := fmt.Sprintf("Port %d %s on %s", change.Port, action, change.Host)
	title := fmt.Sprintf("portwatch: port %s", action)

	body := fmt.Sprintf(`{"token":%q,"user":%q,"title":%q,"message":%q}`,
		p.token, p.user, title, message)

	resp, err := p.client.Post(p.apiURL, "application/json", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushover: request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Status int `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("pushover: failed to decode response: %w", err)
	}
	if result.Status != 1 {
		return fmt.Errorf("pushover: unexpected status %d", result.Status)
	}
	return nil
}
