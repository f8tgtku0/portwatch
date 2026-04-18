package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_Spike(t *testing.T) {
	if !notify.IsSupported("spike") {
		t.Error("expected spike to be a supported channel")
	}
}

func TestSupportedChannels_ContainsSpike(t *testing.T) {
	for _, ch := range notify.SupportedChannels() {
		if ch == "spike" {
			return
		}
	}
	t.Error("spike not found in SupportedChannels()")
}
