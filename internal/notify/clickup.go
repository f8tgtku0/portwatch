package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

const clickupAPIBase = "https://api.clickup.com/api/v2"

type ClickUpNotifier struct {
	token  string
	listID string
	baseURL string
}

func NewClickUpNotifier(token, listID string) *ClickUpNotifier {
	return &ClickUpNotifier{
		token:   token,
		listID:  listID,
		baseURL: clickupAPIBase,
	}
}

func (c *ClickUpNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}

	title := fmt.Sprintf("Port %d %s", change.Port, action)
	description := fmt.Sprintf("Port %d was unexpectedly %s on the monitored host.", change.Port, action)

	body := map[string]interface{}{
		"name":        title,
		"description": description,
		"priority":    2,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("clickup: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/list/%s/task", c.baseURL, c.listID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("clickup: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("clickup: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clickup: unexpected status %d", resp.StatusCode)
	}
	return nil
}
