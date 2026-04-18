package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/state"
)

func TestDebounce_CollapsesFastFlaps(t *testing.T) {
	var mu sync.Mutex
	var received []state.Change

	d := debounce.New(80*time.Millisecond, func(changes []state.Change) {
		mu.Lock()
		received = append(received, changes...)
		mu.Unlock()
	})

	d.Submit(state.Change{Port: 8080, Action: state.Opened})
	d.Submit(state.Change{Port: 8080, Action: state.Closed})
	d.Submit(state.Change{Port: 8080, Action: state.Opened})

	time.Sleep(160 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 flushed change, got %d", len(received))
	}
	if received[0].Action != state.Opened {
		t.Errorf("expected Opened, got %v", received[0].Action)
	}
}

func TestDebounce_IndependentPorts(t *testing.T) {
	var mu sync.Mutex
	var received []state.Change

	d := debounce.New(60*time.Millisecond, func(changes []state.Change) {
		mu.Lock()
		received = append(received, changes...)
		mu.Unlock()
	})

	d.Submit(state.Change{Port: 80, Action: state.Opened})
	d.Submit(state.Change{Port: 443, Action: state.Opened})

	time.Sleep(130 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(received))
	}
}

func TestDebounce_FlushDeliversImmediately(t *testing.T) {
	var mu sync.Mutex
	var received []state.Change

	d := debounce.New(5*time.Second, func(changes []state.Change) {
		mu.Lock()
		received = append(received, changes...)
		mu.Unlock()
	})

	d.Submit(state.Change{Port: 9090, Action: state.Opened})
	d.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 change after Flush, got %d", len(received))
	}
}

func TestDebounce_FlushEmptyIsNoop(t *testing.T) {
	called := false
	d := debounce.New(50*time.Millisecond, func(_ []state.Change) {
		called = true
	})
	d.Flush()
	if called {
		t.Error("flush on empty debouncer should not call downstream")
	}
}
