// Package snapshot captures and persists periodic port-state snapshots.
//
// A Snapshot records the full list of open ports observed during each scan
// cycle and writes them to disk as JSON. This allows portwatch to detect
// long-term drift by comparing the current state against a known-good
// baseline captured at a specific point in time.
//
// Usage:
//
//	snap := snapshot.New("/var/lib/portwatch/snapshot.json")
//	_ = snap.Load()          // restore previous snapshot on startup
//	rec := snapshot.NewRecorder(snap)
//	ports = rec.Apply(ports) // record each scan result transparently
package snapshot
