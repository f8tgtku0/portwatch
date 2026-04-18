package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/state"
)

// NewJiraNotifier creates a Notifier that opens a Jira issue on port changes.
func NewJiraNotifier(baseURL, project, issueType, email, apiToken string) Notifier {
	return &jiraNotifier{
		baseURL:   baseURL,
		project:   project,
		issueType: issueType,
		email:     email,
		apiToken:  apiToken,
		client:    &http.Client{},
	}
}

type jiraNotifier struct {
	baseURL   string
	project   string
	issueType string
	email     string
	apiToken  string
	client    *http.Client
}

func (j *jiraNotifier) Send(change state.Change) error {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	summary := fmt.Sprintf("Port %d unexpectedly %s", change.Port, action)
	body := map[string]any{
		"fields": map[string]any{
			"project":   map[string]string{"key": j.project},
			"summary":   summary,
			"issuetype": map[string]string{"name": j.issueType},
			"description": map[string]any{
				"type":    "doc",
				"version": 1,
				"content": []map[string]any{
					{"type": "paragraph", "content": []map[string]any{
						{"type": "text", "text": summary},
					}},
				},
			},
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, j.baseURL+"/rest/api/3/issue", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(j.email, j.apiToken)
	resp, err := j.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
