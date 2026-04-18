package notify

import (
	"fmt"
)

func init() {
	registerBuilder("jira", buildJira)
}

func buildJira(cfg map[string]string) (Notifier, error) {
	baseURL, ok := cfg["url"]
	if !ok || baseURL == "" {
		return nil, fmt.Errorf("jira: missing required field 'url'")
	}
	project, ok := cfg["project"]
	if !ok || project == "" {
		return nil, fmt.Errorf("jira: missing required field 'project'")
	}
	issueType := cfg["issue_type"]
	if issueType == "" {
		issueType = "Bug"
	}
	email, ok := cfg["email"]
	if !ok || email == "" {
		return nil, fmt.Errorf("jira: missing required field 'email'")
	}
	token, ok := cfg["api_token"]
	if !ok || token == "" {
		return nil, fmt.Errorf("jira: missing required field 'api_token'")
	}
	return NewJiraNotifier(baseURL, project, issueType, email, token), nil
}
