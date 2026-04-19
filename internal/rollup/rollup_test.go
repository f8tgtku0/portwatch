package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/state"
)

func opened(port int) state.Change {
	return state.Change{Port: port, Action: state.Opened}
}

func TestAdd_BatchesChanges(t *testing.T) {
	var got []state.Change
	r := rollup.New(50*time.Millisecond, func(ch []state.Change) {
		got = append(got, ch...)
	})

	r.Add([]state.Change{opened(80)})
	r.Add([]state.Change{opened(443)})

	time.Sleep(100 * time.Millisecond)

	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
}

func TestAdd_EmptyChangesIgnored(t *testing.T) {
	called := false
	r := rollup.New(20*time.Millisecond, func(_ []state.Change) { called = true })
	r.Add(nil)
	time.Sleep(50 * time.Millisecond)
	if called {
		t.Fatal("handler should not be called for empty input")
	}
}

func TestFlush_DeliversImmediately(t *testing.T) {
	var got []state.Change
	r := rollup.New(5*time.Second, func(ch []state.Change) {
		got = ch
	})

	r.Add([]state.Change{opened(22)})
	r.Flush()

	if len(got) != 1 {
		t.Fatalf("expected 1 change after flush, got %d", len(got))
	}
}

func TestFlush_EmptyIsNoop(t *testing.T) {
	called := false
	r := rollup.New(20*time.Millisecond, func(_ []state.Change) { called = true })
	r.Flush()
	if called {
		t.Fatal("flush on empty batch should not invoke handler")
	}
}

func TestMiddleware_Apply_AndFlush(t *testing.T) {
	var got []state.Change
	m := rollup.NewMiddleware(5*time.Second, func(ch []state.Change) { got = ch })
	m.Apply([]state.Change{opened(8080)})
	m.Flush()
	if len(got) != 1 || got[0].Port != 8080 {
		t.Fatalf("unexpected changes: %v", got)
	}
}
