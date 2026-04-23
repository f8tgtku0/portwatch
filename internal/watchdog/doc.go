// Package watchdog provides a self-monitoring heartbeat component for portwatch.
//
// The Watchdog expects periodic Beat() calls from the scan loop. If no beat is
// received within the configured timeout window, it writes a warning to the
// configured writer (defaulting to stderr).
//
// Typical usage:
//
//	wd := watchdog.New(30*time.Second, os.Stderr)
//	go wd.Start()
//	defer wd.Stop()
//
//	// inside scan loop:
//	wd.Beat()
package watchdog
