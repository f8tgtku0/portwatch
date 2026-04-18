package notify

// SupportedChannels lists all valid notification channel identifiers.
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
	"sns",
	"googlechat",
	"desktop",
	"twilio",
	"signalwire",
	"pushbullet",
}

// IsSupported returns true if the given channel name is recognised.
func IsSupported(name string) bool {
	for _, c := range SupportedChannels {
		if c == name {
			return true
		}
	}
	return false
}
