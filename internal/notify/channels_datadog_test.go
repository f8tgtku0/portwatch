package notify_test

import (
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestIsSupported_Datadog(t *testing.T) {
	if !notify.IsSupported("datadog") {
		t.Error("expected datadog to be a supported channel")
	}
}

func TestSupportedChannels_ContainsDatadog(t *testing.T) {
	for _, ch := range notify.SupportedChannels() {
		if ch == "datadog" {
			return
		}
	}
	t.Error("datadog not found in SupportedChannels")
}
