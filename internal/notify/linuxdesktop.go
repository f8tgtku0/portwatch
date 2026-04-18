package notify

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/user/portwatch/internal/state"
)

// DesktopNotifier sends desktop notifications via libnotify (notify-send) on
// Linux or osascript on macOS.
type DesktopNotifier struct {
	appName string
}

// NewDesktopNotifier returns a DesktopNotifier. appName is shown as the
// notification title prefix.
func NewDesktopNotifier(appName string) *DesktopNotifier {
	if appName == "" {
		appName = "portwatch"
	}
	return &DesktopNotifier{appName: appName}
}

// Send dispatches a desktop notification for the given port change.
func (d *DesktopNotifier) Send(change state.Change) error {
	title, body := d.format(change)
	switch runtime.GOOS {
	case "linux":
		return exec.Command("notify-send", title, body).Run()
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, body, title)
		return exec.Command("osascript", "-e", script).Run()
	default:
		return fmt.Errorf("desktop notifications not supported on %s", runtime.GOOS)
	}
}

func (d *DesktopNotifier) format(change state.Change) (title, body string) {
	action := "opened"
	if !change.Opened {
		action = "closed"
	}
	title = fmt.Sprintf("%s – port %s", d.appName, action)
	body = fmt.Sprintf("Port %d (%s) was %s", change.Port.Number, change.Port.Proto, action)
	return
}
