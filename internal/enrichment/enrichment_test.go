package enrichment_test

import (
	"testing"

	"github.com/user/portwatch/internal/enrichment"
)

func TestLookup_WellKnownPort(t *testing.T) {
	e := enrichment.New(nil)
	en := e.Lookup(22)
	if en.Service != "ssh" {
		t.Errorf("expected ssh, got %s", en.Service)
	}
	if en.Proto != "tcp" {
		t.Errorf("expected tcp, got %s", en.Proto)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	e := enrichment.New(nil)
	en := e.Lookup(19999)
	if en.Service != "unknown" {
		t.Errorf("expected unknown, got %s", en.Service)
	}
	if en.Port != 19999 {
		t.Errorf("expected port 19999, got %d", en.Port)
	}
}

func TestLookup_CustomOverridesBuiltin(t *testing.T) {
	custom := map[int]enrichment.Entry{
		80: {Port: 80, Service: "my-app", Proto: "tcp"},
	}
	e := enrichment.New(custom)
	en := e.Lookup(80)
	if en.Service != "my-app" {
		t.Errorf("expected my-app, got %s", en.Service)
	}
}

func TestLookup_CustomNewPort(t *testing.T) {
	custom := map[int]enrichment.Entry{
		9000: {Port: 9000, Service: "myservice", Proto: "tcp"},
	}
	e := enrichment.New(custom)
	en := e.Lookup(9000)
	if en.Service != "myservice" {
		t.Errorf("expected myservice, got %s", en.Service)
	}
}

func TestLabel_WellKnown(t *testing.T) {
	e := enrichment.New(nil)
	label := e.Label(443)
	if label != "443/https" {
		t.Errorf("expected 443/https, got %s", label)
	}
}

func TestLabel_Unknown(t *testing.T) {
	e := enrichment.New(nil)
	label := e.Label(54321)
	if label != "54321" {
		t.Errorf("expected 54321, got %s", label)
	}
}

func TestNew_NilCustom_DoesNotPanic(t *testing.T) {
	e := enrichment.New(nil)
	_ = e.Lookup(80)
}
