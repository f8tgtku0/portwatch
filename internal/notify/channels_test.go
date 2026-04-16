package notify

import "testing"

func TestIsSupported_KnownChannels(t *testing.T) {
	for _, ch := range []string{"log", "slack", "discord", "teams", "pagerduty", "opsgenie", "victorops", "webhook", "email"} {
		if !IsSupported(ch) {
			t.Errorf("expected %q to be supported", ch)
		}
	}
}

func TestIsSupported_UnknownChannel(t *testing.T) {
	if IsSupported("unknown") {
		t.Error("expected 'unknown' to not be supported")
	}
	if IsSupported("") {
		t.Error("expected empty string to not be supported")
	}
}

func TestSupportedChannels_ContainsVictorOps(t *testing.T) {
	found := false
	for _, c := range SupportedChannels {
		if c == "victorops" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected SupportedChannels to contain 'victorops'")
	}
}
