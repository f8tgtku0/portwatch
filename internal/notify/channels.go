package notify

// SupportedChannels lists all available notification channel names.
var SupportedChannels = []string{
	"log",
	"webhook",
	"email",
	"slack",
	"discord",
	"teams",
	"pagerduty",
	"opsgenie",
	"victorops",
}

// IsSupported returns true if the given channel name is supported.
func IsSupported(channel string) bool {
	for _, c := range SupportedChannels {
		if c == channel {
			return true
		}
	}
	return false
}
