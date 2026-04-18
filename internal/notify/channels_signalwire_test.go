package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_SignalWire(t *testing.T) {
	if !notify.IsSupported("signalwire") {
		t.Error("expected signalwire to be a supported channel")
	}
}

func TestSupportedChannels_ContainsSignalWire(t *testing.T) {
	found := false
	for _, ch := range notify.SupportedChannels() {
		if ch == "signalwire" {
			found = true
			break
		}
	}
	if !found {
		t.Error("signalwire missing from SupportedChannels")
	}
}
