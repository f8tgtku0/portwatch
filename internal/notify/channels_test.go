package notify

import "testing"

func TestIsSupported_KnownChannels(t *testing.T) {
	for _, ch := range []string{"log", "webhook", "slack", "email", "discord"} {
		if !IsSupported(ch) {
			t.Errorf("expected %q to be supported", ch)
		}
	}
}

func TestIsSupported_UnknownChannel(t *testing.T) {
	if IsSupported("unknown-channel") {
		t.Error("expected unknown-channel to not be supported")
	}
}

func TestSupportedChannels_ContainsVictorOps(t *testing.T) {
	if !IsSupported("victorops") {
		t.Error("expected victorops to be supported")
	}
}

func TestSupportedChannels_ContainsSMS(t *testing.T) {
	if !IsSupported("sms") {
		t.Error("expected sms to be supported")
	}
}

func TestSupportedChannels_Length(t *testing.T) {
	if len(SupportedChannels) < 14 {
		t.Errorf("expected at least 14 channels, got %d", len(SupportedChannels))
	}
}
