package routing_test

import (
	"testing"

	"github.com/user/portwatch/internal/routing"
	"github.com/user/portwatch/internal/state"
)

func TestLookup_WellKnownPort(t *testing.T) {
	r := routing.New(nil)
	rt := r.Lookup(22)
	if rt.Note != "SSH" {
		t.Errorf("expected SSH, got %q", rt.Note)
	}
	if rt.Direction != routing.Inbound {
		t.Errorf("expected Inbound, got %q", rt.Direction)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	r := routing.New(nil)
	rt := r.Lookup(9999)
	if rt.Note != "" {
		t.Errorf("expected empty note for unknown port, got %q", rt.Note)
	}
	if rt.Port != 9999 {
		t.Errorf("expected port 9999, got %d", rt.Port)
	}
}

func TestLookup_OverridesTakesPriority(t *testing.T) {
	overrides := map[int]routing.Route{
		80: {Port: 80, Direction: routing.Outbound, Protocol: "tcp", Note: "custom"},
	}
	r := routing.New(overrides)
	rt := r.Lookup(80)
	if rt.Note != "custom" {
		t.Errorf("expected custom override, got %q", rt.Note)
	}
	if rt.Direction != routing.Outbound {
		t.Errorf("expected Outbound from override, got %q", rt.Direction)
	}
}

func TestAnnotate_MapsAllPorts(t *testing.T) {
	r := routing.New(nil)
	ports := []state.Port{
		{Number: 22},
		{Number: 443},
		{Number: 7777},
	}
	annotated := r.Annotate(ports)
	if len(annotated) != 3 {
		t.Fatalf("expected 3 annotations, got %d", len(annotated))
	}
	if annotated[22].Note != "SSH" {
		t.Errorf("expected SSH for port 22")
	}
	if annotated[443].Note != "HTTPS" {
		t.Errorf("expected HTTPS for port 443")
	}
	if annotated[7777].Note != "" {
		t.Errorf("expected empty note for unknown port 7777")
	}
}

func TestAnnotate_EmptyPorts(t *testing.T) {
	r := routing.New(nil)
	annotated := r.Annotate([]state.Port{})
	if len(annotated) != 0 {
		t.Errorf("expected empty map, got %d entries", len(annotated))
	}
}
