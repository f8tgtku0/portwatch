package watchdog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestBeat_UpdatesLastBeat(t *testing.T) {
	wd := watchdog.New(5*time.Second, nil)
	before := wd.LastBeat()
	time.Sleep(5 * time.Millisecond)
	wd.Beat()
	after := wd.LastBeat()
	if !after.After(before) {
		t.Errorf("expected LastBeat to advance after Beat(); got before=%v after=%v", before, after)
	}
}

func TestStop_HaltsWatchdog(t *testing.T) {
	var buf bytes.Buffer
	wd := watchdog.New(100*time.Millisecond, &buf)
	wd.Beat()

	done := make(chan struct{})
	go func() {
		wd.Start()
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	wd.Stop()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watchdog did not stop within timeout")
	}
}

func TestStart_WritesWarningWhenStalled(t *testing.T) {
	var buf bytes.Buffer
	timeout := 60 * time.Millisecond
	wd := watchdog.New(timeout, &buf)

	go wd.Start()
	// Do NOT call Beat — let it stall.
	time.Sleep(200 * time.Millisecond)
	wd.Stop()

	output := buf.String()
	if !strings.Contains(output, "stalled") {
		t.Errorf("expected stall warning in output, got: %q", output)
	}
}

func TestStart_NoWarningWhenBeating(t *testing.T) {
	var buf bytes.Buffer
	timeout := 100 * time.Millisecond
	wd := watchdog.New(timeout, &buf)

	go wd.Start()

	// Continuously beat to keep watchdog happy.
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.After(20 * time.Millisecond):
				wd.Beat()
			}
		}
	}()

	time.Sleep(250 * time.Millisecond)
	close(stop)
	wd.Stop()

	if buf.Len() > 0 {
		t.Errorf("unexpected warning output: %q", buf.String())
	}
}

func TestNew_NilWriter_DefaultsToStderr(t *testing.T) {
	// Should not panic with nil writer.
	wd := watchdog.New(1*time.Second, nil)
	wd.Beat() // must not panic
}
