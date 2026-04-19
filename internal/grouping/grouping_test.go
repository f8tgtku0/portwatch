package grouping_test

import (
	"testing"

	"github.com/user/portwatch/internal/grouping"
	"github.com/user/portwatch/internal/state"
)

func TestLabel_MatchesKnownPort(t *testing.T) {
	g := grouping.New([]grouping.Group{
		{Name: "web", Ports: []int{80, 443}},
		{Name: "db", Ports: []int{5432, 3306}},
	})
	if got := g.Label(80); got != "web" {
		t.Fatalf("expected web, got %s", got)
	}
	if got := g.Label(5432); got != "db" {
		t.Fatalf("expected db, got %s", got)
	}
}

func TestLabel_UnknownPort(t *testing.T) {
	g := grouping.New(nil)
	if got := g.Label(9999); got != "unknown" {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestAnnotate_MapsChanges(t *testing.T) {
	g := grouping.New([]grouping.Group{
		{Name: "web", Ports: []int{80}},
	})
	changes := []state.Change{
		{Port: 80, Action: state.Opened},
		{Port: 9999, Action: state.Opened},
	}
	anno := g.Annotate(changes)
	if anno[80] != "web" {
		t.Fatalf("expected web for 80")
	}
	if anno[9999] != "unknown" {
		t.Fatalf("expected unknown for 9999")
	}
}

func TestAdd_AppendsGroup(t *testing.T) {
	g := grouping.New(nil)
	g.Add(grouping.Group{Name: "cache", Ports: []int{6379}})
	if got := g.Label(6379); got != "cache" {
		t.Fatalf("expected cache, got %s", got)
	}
}
