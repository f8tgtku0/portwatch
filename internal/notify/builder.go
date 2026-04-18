package notify

import (
	"fmt"
	"io"

	"github.com/user/portwatch/internal/config"
)

// Build constructs a Notifier from the provided config, writing log output to w.
func Build(cfg config.Config, w io.Writer) (Notifier, error) {
	var notifiers []Notifier
	for _, ch := range cfg.Channels {
		switch ch.Name {
		case "log":
			notifiers = append(notifiers, NewLogNotifier(w))
		case "webhook":
			notifiers = append(notifiers, NewWebhookNotifier(ch.URL))
		case "slack":
			notifiers = append(notifiers, NewSlackNotifier(ch.URL))
		case "discord":
			notifiers = append(notifiers, NewDiscordNotifier(ch.URL))
		case "teams":
			notifiers = append(notifiers, NewTeamsNotifier(ch.URL))
		case "googlechat":
			notifiers = append(notifiers, NewGoogleChatNotifier(ch.URL))
		case "mattermost":
			notifiers = append(notifiers, NewMattermostNotifier(ch.URL))
		case "rocketchat":
			notifiers = append(notifiers, NewRocketChatNotifier(ch.URL))
		case "lark":
			notifiers = append(notifiers, NewLarkNotifier(ch.URL))
		case "ntfy":
			notifiers = append(notifiers, NewNtfyNotifier(ch.URL))
		case "gotify":
			notifiers = append(notifiers, NewGotifyNotifier(ch.URL, ch.Token))
		case "telegram":
			notifiers = append(notifiers, NewTelegramNotifier(ch.Token, ch.ChatID))
		case "pushover":
			notifiers = append(notifiers, NewPushoverNotifier(ch.Token, ch.UserKey))
		case "pushbullet":
			notifiers = append(notifiers, NewPushbulletNotifier(ch.Token))
		case "pagerduty":
			notifiers = append(notifiers, NewPagerDutyNotifier(ch.Token))
		case "opsgenie":
			notifiers = append(notifiers, NewOpsGenieNotifier(ch.Token))
		case "victorops":
			notifiers = append(notifiers, NewVictorOpsNotifier(ch.URL))
		case "desktop":
			notifiers = append(notifiers, NewDesktopNotifier(""))
		default:
			return nil, fmt.Errorf("notify: unsupported channel %q", ch.Name)
		}
	}
	if len(notifiers) == 0 {
		return NewLogNotifier(w), nil
	}
	if len(notifiers) == 1 {
		return notifiers[0], nil
	}
	return NewMulti(notifiers...), nil
}
