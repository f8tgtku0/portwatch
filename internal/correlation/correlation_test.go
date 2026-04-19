package correlation_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlation"
	"github.com/user/portwatch/internal/state"
)

func opened(port int) state.Change { return state.Change{Port: port, Action: state.Opened} }
func closed(port int) state.Change { return state.Change{Port: port, Action: state.Closed} }

func TestAdd_BuffersAndEmitsAfterWindow(t *testing.T) {
	c := correlation.New(50 * time.Millisecond)
	c.Add(opened(80))
	c.Add(closed(80))

	select {
	case ev := <-c.Events():
		t.Fatal("unexpected early emit", ev)
	case <-time.After(10 * time.Millisecond):
	}

	time.Sleep(80 * time.Millisecond)

	select {
	case ev := <-c.Events():
		if len(ev.Changes) != 2 {
			t.Fatalf("expected 2 changes, got %d", len(ev.Changes))
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected event not received")
	}
}

func TestFlush_EmitsImmediately(t *testing.T) {
	c := correlation.New(10 * time.Second)
	c.Add(opened(443))
	c.Flush()

	select {
	case ev := <-c.Events():
		if len(ev.Changes) != 1 {
			t.Fatalf("expected 1 change, got %d", len(ev.Changes))
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("expected flushed event")
	}
}

func TestFlush_EmptyIsNoop(t *testing.T) {
	c := correlation.New(10 * time.Second)
	c.Flush()

	select {
	case ev := <-c.Events():
		t.Fatal("unexpected event", ev)
	case <-time.After(20 * time.Millisecond):
	}
}

func TestAdd_IndependentPorts(t *testing.T) {
	c := correlation.New(50 * time.Millisecond)
	c.Add(opened(22))
	c.Add(opened(3306))
	c.Flush()

	count := 0
	for count < 2 {
		select {
		case <-c.Events():
			count++
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("only received %d events, expected 2", count)
		}
	}
}
