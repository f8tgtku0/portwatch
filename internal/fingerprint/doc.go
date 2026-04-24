// Package fingerprint provides deterministic hashing of open-port snapshots.
//
// A Fingerprinter converts a slice of state.Port values into a stable
// hex-encoded SHA-256 digest.  Because ports are sorted before hashing,
// the same logical set always produces the same fingerprint regardless of
// the order in which ports were discovered.
//
// Typical use:
//
//	fp := fingerprint.New()
//	if fp.Changed(previousPorts, currentPorts) {
//		// run diff + alert pipeline
//	}
package fingerprint
