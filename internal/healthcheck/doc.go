// Package healthcheck exposes a lightweight HTTP /healthz endpoint that
// reports the liveness of the portwatch daemon.
//
// Usage:
//
//	hs := healthcheck.New(":9090")
//	go hs.ListenAndServe()
//
//	// After each scan cycle:
//	hs.RecordScan()
//
// The endpoint returns JSON:
//
//	{
//	  "ok": true,
//	  "uptime": "2m30s",
//	  "last_scan": "2024-01-15T10:00:00Z",
//	  "scan_count": 42
//	}
package healthcheck
