package notify

// SupportedChannels lists all available notification channel identifiers.
var SupportedChannels = []string{
	"log",
	"webhook",
	"email",
	"slack",
	"pagerduty",
	"opsgenie",
	"discord",
	"teams",
	"victorops",
	"telegram",
	"gotify",
	"ntfy",
	"matrix",
	"sms",
	"pushover",
	"mattermost",
	"signalr",
}

// IsSupported returns true if the given channel name is recognised.
func IsSupported(name string) bool {
	for _, ch := range SupportedChannels {
		if ch == name {
			return true
		}
	}
	return false
}
