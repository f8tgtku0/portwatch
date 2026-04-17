package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewLarkNotifier sends messages to a Lark (Feishu) webhook.
func NewLarkNotifier(webhookURL string) Notifier {
	return &larkNotifier{webhookURL: webhookURL, client: &http.Client{}}
}

type larkNotifier struct {
	webhookURL string
	client     *http.Client
}

func (l *larkNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	text := fmt.Sprintf("Port %d (%s) was %s.", change.Port, change.Proto, action)

	body, _ := json.Marshal(map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": text},
	})

	resp, err := l.client.Post(l.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark: unexpected status %d", resp.StatusCode)
	}
	return nil
}
