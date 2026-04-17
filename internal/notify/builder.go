package notify

import (
	"fmt"
	"io"

	"github.com/user/portwatch/internal/config"
)

// Build constructs a Notifier from a config.NotifyChannel entry.
// It returns an error if the channel type is unsupported or misconfigured.
func Build(ch config.NotifyChannel, w io.Writer) (Notifier, error) {
	switch ch.Type {
	case "log":
		return NewLogNotifier(w), nil
	case "webhook":
		return NewWebhookNotifier(ch.URL), nil
	case "slack":
		return NewSlackNotifier(ch.URL), nil
	case "discord":
		return NewDiscordNotifier(ch.URL), nil
	case "teams":
		return NewTeamsNotifier(ch.URL), nil
	case "pagerduty":
		return NewPagerDutyNotifier(ch.URL), nil
	case "opsgenie":
		return NewOpsGenieNotifier(ch.URL), nil
	case "victorops":
		return NewVictorOpsNotifier(ch.URL), nil
	case "telegram":
		return NewTelegramNotifier(ch.URL), nil
	case "gotify":
		return NewGotifyNotifier(ch.URL), nil
	case "ntfy":
		return NewNtfyNotifier(ch.URL), nil
	case "matrix":
		return NewMatrixNotifier(ch.URL), nil
	case "mattermost":
		return NewMattermostNotifier(ch.URL), nil
	case "rocketchat":
		return NewRocketChatNotifier(ch.URL), nil
	case "zulip":
		return NewZulipNotifier(ch.URL), nil
	case "lark":
		return NewLarkNotifier(ch.URL), nil
	case "signalr":
		return NewSignalRNotifier(ch.URL), nil
	case "sns":
		return NewSNSNotifier(ch.URL)
	case "email":
		return NewEmailNotifier(ch.Host, ch.From, ch.To), nil
	case "sms":
		return NewSMSNotifier(ch.URL, ch.From, ch.To), nil
	case "pushover":
		return NewPushoverNotifier(ch.URL), nil
	default:
		return nil, fmt.Errorf("notify: unsupported channel type %q", ch.Type)
	}
}
