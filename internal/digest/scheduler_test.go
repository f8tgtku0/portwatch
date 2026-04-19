package digest_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/digest"
)

func TestScheduler_FlushesOnInterval(t *testing.T) {
	var buf bytes.Buffer
	d := digest.New(&buf)
	d.Record(openedPorts(80), nil)

	s := digest.NewScheduler(d, 30*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	if !strings.Contains(buf.String(), "opened") {
		t.Errorf("expected flush output, got: %s", buf.String())
	}
}

func TestScheduler_FlushesOnCancel(t *testing.T) {
	var buf bytes.Buffer
	d := digest.New(&buf)
	d.Record(openedPorts(443), nil)

	s := digest.NewScheduler(d, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	s.Run(ctx)

	if !strings.Contains(buf.String(), "443") {
		t.Errorf("expected port 443 in flush output on cancel, got: %s", buf.String())
	}
}
