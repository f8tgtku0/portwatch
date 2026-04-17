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
	"rocketchat",
	"zulip",
	"lark",
}

// IsSupported returns true if the given channel name is a known notifier.
func IsSupported(channel string) bool {
	for _, c := range SupportedChannels {
		if c == channel {
			return true
		}
	}
	return false
}
