// Package ratelimit implements cooldown-based suppression of repeated port
// change alerts. It prevents notification floods when a port flaps or when
// the monitor restarts and re-detects the same open ports.
//
// Usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//	for _, change := range changes {
//		key := ratelimit.ChangeKey(change)
//		if limiter.Allow(key) {
//			alert.Notify(change)
//		}
//	}
package ratelimit
