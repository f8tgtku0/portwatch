// Package baseline provides a trusted-port registry for portwatch.
//
// Ports added to the baseline are considered expected and will not trigger
// alerts when they appear in scan diffs. The baseline is persisted to a JSON
// file so it survives daemon restarts.
//
// Typical usage:
//
//	b := baseline.New("/var/lib/portwatch/baseline.json")
//	_ = b.Load()
//	mw := baseline.NewMiddleware(b)
//	filtered := mw.Apply(changes)
package baseline
