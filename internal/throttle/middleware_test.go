package throttle_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/throttle"
)

type countNotifier struct{ count int }

func (c *countNotifier) Send(_ context.Context, _ notify.Message) error {
	c.count++
	return nil
}

func openedMsg(port int) notify.Message {
	return notify.Message{Change: state.Change{Port: port, Opened: true}}
}

func TestMiddleware_AllowsFirstSend(t *testing.T) {
	n := &countNotifier{}
	mw := throttle.NewMiddleware(n, throttle.New(time.Second, 1))
	_ = mw.Send(context.Background(), openedMsg(8080))
	if n.count != 1 {
		t.Fatalf("expected 1 send, got %d", n.count)
	}
}

func TestMiddleware_ThrottlesExcessiveSends(t *testing.T) {
	n := &countNotifier{}
	mw := throttle.NewMiddleware(n, throttle.New(time.Second, 1))
	for i := 0; i < 5; i++ {
		_ = mw.Send(context.Background(), openedMsg(9090))
	}
	if n.count != 1 {
		t.Fatalf("expected 1 send due to throttle, got %d", n.count)
	}
}

func TestMiddleware_AllowsAfterWindowExpires(t *testing.T) {
	n := &countNotifier{}
	mw := throttle.NewMiddleware(n, throttle.New(50*time.Millisecond, 1))
	_ = mw.Send(context.Background(), openedMsg(3000))
	time.Sleep(60 * time.Millisecond)
	_ = mw.Send(context.Background(), openedMsg(3000))
	if n.count != 2 {
		t.Fatalf("expected 2 sends after window reset, got %d", n.count)
	}
}
