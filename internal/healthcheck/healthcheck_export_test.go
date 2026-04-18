package healthcheck

// CurrentStatus returns a snapshot of the current health status.
// Exported only for testing.
func (s *Server) CurrentStatus() Status {
	uptime := clock().Sub(s.start).Round(time.Second).String()
	st := Status{
		OK:        true,
		Uptime:    uptime,
		ScanCount: s.scanCount.Load(),
	}
	if ls, ok := s.lastScan.Load().(time.Time); ok && !ls.IsZero() {
		st.LastScan = ls
	}
	return st
}
