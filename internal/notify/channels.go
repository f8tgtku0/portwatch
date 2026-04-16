// Package notify provides notifier implementations for various channels.
// This file documents and registers the available notifier constructors.
package notify

// Channel represents a supported notification channel name.
type Channel string

const (
	ChannelLog       Channel = "log"
	ChannelWebhook   Channel = "webhook"
	ChannelSlack     Channel = "slack"
	ChannelDiscord   Channel = "discord"
	ChannelTeams     Channel = "teams"
	ChannelEmail     Channel = "email"
	ChannelPagerDuty Channel = "pagerduty"
	ChannelOpsGenie  Channel = "opsgenie"
)

// SupportedChannels returns a list of all supported notification channel names.
func SupportedChannels() []Channel {
	return []Channel{
		ChannelLog,
		ChannelWebhook,
		ChannelSlack,
		ChannelDiscord,
		ChannelTeams,
		ChannelEmail,
		ChannelPagerDuty,
		ChannelOpsGenie,
	}
}

// IsSupported returns true if the given channel name is recognised.
func IsSupported(c Channel) bool {
	for _, s := range SupportedChannels() {
		if s == c {
			return true
		}
	}
	return false
}
