package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_ClickUp(t *testing.T) {
	if !notify.IsSupported("clickup") {
		t.Error("expected clickup to be a supported channel")
	}
}

func TestSupportedChannels_ContainsClickUp(t *testing.T) {
	channels := notify.SupportedChannels()
	for _, ch := range channels {
		if ch == "clickup" {
			return
		}
	}
	t.Error("expected clickup in supported channels list")
}
