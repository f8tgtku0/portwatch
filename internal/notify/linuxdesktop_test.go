package notify

import (
	"runtime"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestNewDesktopNotifier_DefaultAppName(t *testing.T) {
	d := NewDesktopNotifier("")
	if d.appName != "portwatch" {
		t.Errorf("expected default appName 'portwatch', got %q", d.appName)
	}
}

func TestNewDesktopNotifier_CustomAppName(t *testing.T) {
	d := NewDesktopNotifier("myapp")
	if d.appName != "myapp" {
		t.Errorf("expected appName 'myapp', got %q", d.appName)
	}
}

func TestDesktopNotifier_Format_Opened(t *testing.T) {
	d := NewDesktopNotifier("portwatch")
	change := state.Change{Port: state.Port{Number: 8080, Proto: "tcp"}, Opened: true}
	title, body := d.format(change)
	if title != "portwatch – port opened" {
		t.Errorf("unexpected title: %q", title)
	}
	if body != "Port 8080 (tcp) was opened" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestDesktopNotifier_Format_Closed(t *testing.T) {
	d := NewDesktopNotifier("portwatch")
	change := state.Change{Port: state.Port{Number: 22, Proto: "tcp"}, Opened: false}
	title, body := d.format(change)
	if title != "portwatch – port closed" {
		t.Errorf("unexpected title: %q", title)
	}
	if body != "Port 22 (tcp) was closed" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestDesktopNotifier_Send_UnsupportedOS(t *testing.T) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		t.Skip("skipping unsupported-OS test on linux/darwin")
	}
	d := NewDesktopNotifier("portwatch")
	change := state.Change{Port: state.Port{Number: 80, Proto: "tcp"}, Opened: true}
	err := d.Send(change)
	if err == nil {
		t.Error("expected error for unsupported OS, got nil")
	}
}

func TestNewDesktopNotifier_ImplementsNotifier(t *testing.T) {
	var _ Notifier = NewDesktopNotifier("portwatch")
}
