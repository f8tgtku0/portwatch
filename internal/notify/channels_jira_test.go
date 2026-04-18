package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_Jira(t *testing.T) {
	if !notify.IsSupported("jira") {
		t.Error("expected jira to be a supported channel")
	}
}

func TestSupportedChannels_ContainsJira(t *testing.T) {
	for _, ch := range notify.SupportedChannels() {
		if ch == "jira" {
			return
		}
	}
	t.Error("jira not found in SupportedChannels")
}
