package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/history"
)

const telegramAPIBase = "https://api.telegram.org"

type TelegramNotifier struct {
	token  string
	chatID string
	baseURL string
	client *http.Client
}

func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:   token,
		chatID:  chatID,
		baseURL: telegramAPIBase,
		client:  &http.Client{},
	}
}

func (t *TelegramNotifier) Send(e history.Entry) error {
	action := "closed"
	if e.Change == "opened" {
		action = "opened ✅"
	} else {
		action = "closed ⚠️"
	}

	text := fmt.Sprintf("*portwatch alert*\nPort `%d` %s at `%s`", e.Port, action, e.Time.Format("2006-01-02 15:04:05"))

	payload := map[string]string{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", t.baseURL, t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status: %d", resp.StatusCode)
	}
	return nil
}
