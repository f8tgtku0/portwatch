// Package healthcheck provides a simple HTTP health endpoint for portwatch.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health state of the daemon.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	LastScan  time.Time `json:"last_scan,omitempty"`
	ScanCount int64     `json:"scan_count"`
}

// Server exposes a /healthz HTTP endpoint.
type Server struct {
	addr      string
	start     time.Time
	scanCount atomic.Int64
	lastScan  atomic.Value // stores time.Time
}

// New creates a new health check server listening on addr (e.g. ":9090").
func New(addr string) *Server {
	return &Server{
		addr:  addr,
		start: time.Now(),
	}
}

// RecordScan updates the last scan timestamp and increments the counter.
func (s *Server) RecordScan() {
	s.scanCount.Add(1)
	s.lastScan.Store(time.Now())
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.start).Round(time.Second).String()

	st := Status{
		OK:        true,
		Uptime:    uptime,
		ScanCount: s.scanCount.Load(),
	}

	if ls, ok := s.lastScan.Load().(time.Time); ok && !ls.IsZero() {
		st.LastScan = ls
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(st)
}
