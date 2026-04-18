package notify

import "testing"

func TestIsSupported_Pushbullet(t *testing.T) {
	if !IsSupported("pushbullet") {
		t.Error("expected pushbullet to be a supported channel")
	}
}

func TestSupportedChannels_ContainsPushbullet(t *testing.T) {
	found := false
	for _, c := range SupportedChannels {
		if c == "pushbullet" {
			found = true
			break
		}
	}
	if !found {
		t.Error("SupportedChannels does not contain pushbullet")
	}
}
