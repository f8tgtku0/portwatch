package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_BearyChat(t *testing.T) {
	if !notify.IsSupported("bearychat") {
		t.Error("expected bearychat to be a supported channel")
	}
}

func TestSupportedChannels_ContainsBearyChat(t *testing.T) {
	for _, ch := range notify.SupportedChannels() {
		if ch == "bearychat" {
			return
		}
	}
	t.Error("bearychat not found in SupportedChannels")
}
