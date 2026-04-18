package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/state"
)

func TestMiddleware_Apply_DebouncesThenForwards(t *testing.T) {
	var mu sync.Mutex
	var got []state.Change

	mw := debounce.NewMiddleware(70*time.Millisecond, func(changes []state.Change) {
		mu.Lock()
		got = append(got, changes...)
		mu.Unlock()
	})

	mw.Apply([]state.Change{
		{Port: 22, Action: state.Opened},
		{Port: 80, Action: state.Opened},
	})

	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 forwarded changes, got %d", len(got))
	}
}

func TestMiddleware_Flush_ForcesSend(t *testing.T) {
	var mu sync.Mutex
	var got []state.Change

	mw := debounce.NewMiddleware(10*time.Second, func(changes []state.Change) {
		mu.Lock()
		got = append(got, changes...)
		mu.Unlock()
	})

	mw.Apply([]state.Change{{Port: 3000, Action: state.Closed}})
	mw.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 change after Flush, got %d", len(got))
	}
}
