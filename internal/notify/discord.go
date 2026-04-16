package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/history"
)

// NewDiscordNotifier sends port change alerts to a Discord webhook.
func NewDiscordNotifier(webhookURL string) Notifier {
	return &discordNotifier{webhookURL: webhookURL, client: &http.Client{Timeout: 5 * time.Second}}
}

type discordNotifier struct {
	webhookURL string
	client     *http.Client
}

type discordPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

func (d *discordNotifier) Send(msg history.Entry) error {
	color := 0x2ecc71 // green for opened
	if msg.Change == "closed" {
		color = 0xe74c3c // red for closed
	}
	payload := discordPayload{
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("Port %s: %d/%s", msg.Change, msg.Port, msg.Proto),
				Description: fmt.Sprintf("Port **%d** (%s) was **%s** at %s", msg.Port, msg.Proto, msg.Change, msg.Timestamp.Format(time.RFC3339)),
				Color:       color,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal: %w", err)
	}
	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
