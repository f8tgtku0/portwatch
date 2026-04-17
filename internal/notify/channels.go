package notify

// SupportedChannels lists all notification channel identifiers recognised by portwatch.
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
}

// IsSupported reports whether the given channel name is supported.
func IsSupported(channel string) bool {
	for _, c := range SupportedChannels {
		if c == channel {
			return true
		}
	}
	return false
}
