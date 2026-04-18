// Package suppress implements maintenance-window suppression for portwatch alerts.
//
// During a suppression window, change notifications are silenced. This is useful
// when planned maintenance would otherwise trigger a flood of alerts.
//
// Usage:
//
//	s := suppress.New()
//	s.Add(start, end, "planned maintenance")
//	if ok, w := s.IsSuppressed(time.Now()); ok {
//		log.Printf("suppressed: %s", w.Reason)
//	}
package suppress
