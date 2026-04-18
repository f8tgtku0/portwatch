package ratelimit_test

import (
	"testing"

	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/state"
)

func TestChangeKey_OpenedPort(t *testing.T) {
	c := state.Change{Port: 8080, Opened: true}
	key := ratelimit.ChangeKey(c)
	if key != "port:8080:opened" {
		t.Fatalf("unexpected key: %s", key)
	}
}

func TestChangeKey_ClosedPort(t *testing.T) {
	c := state.Change{Port: 443, Opened: false}
	key := ratelimit.ChangeKey(c)
	if key != "port:443:closed" {
		t.Fatalf("unexpected key: %s", key)
	}
}

func TestChangeKey_UniquePerPortAndAction(t *testing.T) {
	open := ratelimit.ChangeKey(state.Change{Port: 22, Opened: true})
	closed := ratelimit.ChangeKey(state.Change{Port: 22, Opened: false})
	if open == closed {
		t.Fatal("expected open and closed keys to differ")
	}
}
