package digest

import (
	"context"
	"time"
)

// Scheduler periodically flushes a Digest on a fixed interval.
type Scheduler struct {
	digest   *Digest
	interval time.Duration
}

// NewScheduler creates a Scheduler that flushes d every interval.
func NewScheduler(d *Digest, interval time.Duration) *Scheduler {
	return &Scheduler{digest: d, interval: interval}
}

// Run starts the flush loop and blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.digest.Flush()
		case <-ctx.Done():
			s.digest.Flush()
			return
		}
	}
}
