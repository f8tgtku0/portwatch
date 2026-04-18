package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_Zendesk(t *testing.T) {
	if !notify.IsSupported("zendesk") {
		t.Error("expected zendesk to be a supported channel")
	}
}

func TestSupportedChannels_ContainsZendesk(t *testing.T) {
	for _, ch := range notify.SupportedChannels() {
		if ch == "zendesk" {
			return
		}
	}
	t.Error("zendesk not found in SupportedChannels")
}
