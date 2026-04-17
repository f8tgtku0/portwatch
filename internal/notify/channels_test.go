package notify

import "testing"

func TestIsSupported_KnownChannels(t *testing.T) {
	for _, ch := range []string{"slack", "webhook", "email", "sns"} {
		if !IsSupported(ch) {
			t.Errorf("expected %q to be supported", ch)
		}
	}
}

func TestIsSupported_UnknownChannel(t *testing.T) {
	if IsSupported("carrier-pigeon") {
		t.Error("expected unknown channel to return false")
	}
}

func TestSupportedChannels_ContainsVictorOps(t *testing.T) {
	if !IsSupported("victorops") {
		t.Error("victorops should be supported")
	}
}

func TestSupportedChannels_ContainsSMS(t *testing.T) {
	if !IsSupported("sms") {
		t.Error("sms should be supported")
	}
}

func TestSupportedChannels_Length(t *testing.T) {
	if len(SupportedChannels) < 20 {
		t.Errorf("expected at least 20 channels, got %d", len(SupportedChannels))
	}
}

func TestSupportedChannels_ContainsSNS(t *testing.T) {
	if !IsSupported("sns") {
		t.Error("sns should be in supported channels")
	}
}
