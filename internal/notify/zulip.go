package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewZulipNotifier sends messages to a Zulip stream via the Zulip REST API.
func NewZulipNotifier(baseURL, botEmail, botAPIKey, stream, topic string) Notifier {
	return &zulipNotifier{
		baseURL:   baseURL,
		botEmail:  botEmail,
		botAPIKey: botAPIKey,
		stream:    stream,
		topic:     topic,
		client:    &http.Client{},
	}
}

type zulipNotifier struct {
	baseURL   string
	botEmail  string
	botAPIKey string
	stream    string
	topic     string
	client    *http.Client
}

func (z *zulipNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	content := fmt.Sprintf("Port **%d** (%s) was **%s**.", change.Port, change.Proto, action)

	body, _ := json.Marshal(map[string]string{
		"type":    "stream",
		"to":      z.stream,
		"topic":   z.topic,
		"content": content,
	})

	req, err := http.NewRequest(http.MethodPost, z.baseURL+"/api/v1/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(z.botEmail, z.botAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}
