package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.PortRange.Start != 1 || cfg.PortRange.End != 1024 {
		t.Errorf("unexpected default port range: %+v", cfg.PortRange)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %v", cfg.Interval)
	}
	if cfg.Timeout != 500*time.Millisecond {
		t.Errorf("unexpected default timeout: %v", cfg.Timeout)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := map[string]any{
		"port_range": map[string]any{"start": 80, "end": 443},
		"interval":   float64(60 * time.Second),
		"timeout":    float64(time.Second),
	}
	path := writeTempConfig(t, raw)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.PortRange.Start != 80 || cfg.PortRange.End != 443 {
		t.Errorf("unexpected port range: %+v", cfg.PortRange)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	cfg := config.Default()
	cfg.PortRange.Start = 1000
	cfg.PortRange.End = 500
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for inverted range")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for zero interval")
	}
}

func TestValidate_OutOfBoundsPort(t *testing.T) {
	cfg := config.Default()
	cfg.PortRange.End = 70000
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for port > 65535")
	}
}
