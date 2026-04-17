package notify

// SupportedChannels lists all available notification channel identifiers.
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
	"telegram",
	"gotify",
	"ntfy",
	"matrix",
	"sms",
}

// IsSupported returns true if the given channel name is recognised.
func IsSupported(channel string) bool {
	for _, c := range SupportedChannels {
		if c == channel {
			return true
		}
	}
	return false
}
