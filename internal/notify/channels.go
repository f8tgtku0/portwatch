package notify

// SupportedChannels lists all available notification channel names.
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
