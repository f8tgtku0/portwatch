package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	PortRange  PortRange     `json:"port_range"`
	Interval   time.Duration `json:"interval"`
	AlertFile  string        `json:"alert_file,omitempty"`
	Timeout    time.Duration `json:"timeout"`
}

// PortRange defines the start and end of the port scan range.
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		PortRange: PortRange{Start: 1, End: 1024},
		Interval:  30 * time.Second,
		Timeout:   500 * time.Millisecond,
	}
}

// Load reads and parses a JSON config file from the given path.
// Missing fields fall back to defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are sane.
func (c *Config) Validate() error {
	if c.PortRange.Start < 1 || c.PortRange.Start > 65535 {
		return fmt.Errorf("config: port_range.start must be between 1 and 65535")
	}
	if c.PortRange.End < 1 || c.PortRange.End > 65535 {
		return fmt.Errorf("config: port_range.end must be between 1 and 65535")
	}
	if c.PortRange.Start > c.PortRange.End {
		return fmt.Errorf("config: port_range.start must be <= port_range.end")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("config: interval must be positive")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("config: timeout must be positive")
	}
	return nil
}
