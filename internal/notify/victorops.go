package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// NewVictorOpsNotifier creates a Notifier that sends alerts to VictorOps (Splunk On-Call).
func NewVictorOpsNotifier(webhookURL string) Notifier {
	return &victorOpsNotifier{webhookURL: webhookURL}
}

type victorOpsNotifier struct {
	webhookURL string
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

func (v *victorOpsNotifier) Send(msg alert.Message) error {
	msgType := "INFO"
	if msg.Level == alert.LevelCritical {
		msgType = "CRITICAL"
	} else if msg.Level == alert.LevelWarning {
		msgType = "WARNING"
	}

	payload := victorOpsPayload{
		MessageType:       msgType,
		EntityID:          fmt.Sprintf("portwatch-port-%d", msg.Port),
		EntityDisplayName: fmt.Sprintf("Port %d %s", msg.Port, msg.Change),
		StateMessage:      msg.Text,
		Timestamp:         time.Now().Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}

	resp, err := http.Post(v.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("victorops: unexpected status %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}
	return nil
}
