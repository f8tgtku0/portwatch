package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewMatrixNotifier sends notifications to a Matrix room via the Client-Server API.
func NewMatrixNotifier(homeserver, roomID, accessToken string) Notifier {
	return &matrixNotifier{
		homeserver:  homeserver,
		roomID:      roomID,
		accessToken: accessToken,
		client:      &http.Client{},
	}
}

type matrixNotifier struct {
	homeserver  string
	roomID      string
	accessToken string
	client      *http.Client
}

type matrixMessage struct {
	MsgType string `json:"msgtype"`
	Body    string `json:"body"`
}

func (m *matrixNotifier) Send(msg Message, change state.Change) error {
	payload := matrixMessage{
		MsgType: "m.text",
		Body:    fmt.Sprintf("%s — port %d (%s)", msg.Title, change.Port, change.Proto),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("matrix: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message", m.homeserver, m.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("matrix: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.accessToken)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}
