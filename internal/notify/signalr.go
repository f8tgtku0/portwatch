package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/history"
)

// NewSignalRNotifier sends port change alerts to a SignalR-compatible HTTP endpoint.
func NewSignalRNotifier(endpointURL, hubMethod string) Notifier {
	return &signalRNotifier{url: endpointURL, method: hubMethod}
}

type signalRNotifier struct {
	url    string
	method string
}

type signalRPayload struct {
	Method  string `json:"method"`
	Message string `json:"message"`
	Action  string `json:"action"`
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
}

func (s *signalRNotifier) Send(e history.Entry) error {
	action := "opened"
	if e.Change == history.Closed {
		action = "closed"
	}

	payload := signalRPayload{
		Method:  s.method,
		Message: fmt.Sprintf("Port %d/%s was %s", e.Port, e.Proto, action),
		Action:  action,
		Port:    e.Port,
		Proto:   e.Proto,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalr: marshal: %w", err)
	}

	resp, err := http.Post(s.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status %d", resp.StatusCode)
	}
	return nil
}
