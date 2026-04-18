package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

type linearNotifier struct {
	apiKey  string
	teamID  string
	apiURL  string
	client  *http.Client
}

type linearIssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TeamID      string `json:"teamId"`
}

// NewLinearNotifier creates a Notifier that opens Linear issues on port changes.
func NewLinearNotifier(apiKey, teamID string) Notifier {
	return &linearNotifier{
		apiKey: apiKey,
		teamID: teamID,
		apiURL: "https://api.linear.app/graphql",
		client: &http.Client{},
	}
}

func (n *linearNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	title := fmt.Sprintf("Port %d unexpectedly %s", change.Port, action)
	desc := fmt.Sprintf("portwatch detected port %d was %s on this host.", change.Port, action)

	query := fmt.Sprintf(`mutation { issueCreate(input: {title: %q, description: %q, teamId: %q}) { success } }`,
		title, desc, n.teamID)

	body, _ := json.Marshal(map[string]string{"query": query})
	req, err := http.NewRequest(http.MethodPost, n.apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("linear: unexpected status %d", resp.StatusCode)
	}
	return nil
}
