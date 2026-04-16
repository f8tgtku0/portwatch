package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the full portwatch configuration.
type Config struct {
	PortRange  [2]int   `json:"port_range"`
	Interval   int      `json:"interval_seconds"`
	StateFile  string   `json:"state_file"`
	HistFile   string   `json:"history_file"`
	MaxHistory int      `json:"max_history"`
	LogFile    string   `json:"log_file"`
	WebhookURL string   `json:"webhook_url"`
	SlackURL   string   `json:"slack_url"`
	Email      *Email   `json:"email,omitempty"`
	AlertOn    []string `json:"alert_on"`
}

// Email holds SMTP configuration for email alerts.
type Email struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
	From     string `json:"from"`
	To       string `json:"to"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		PortRange:  [2]int{1, 65535},
		Interval:   60,
		StateFile:  ".portwatch_state.json",
		HistFile:   ".portwatch_history.json",
		MaxHistory: 500,
		AlertOn:    []string{"opened", "closed"},
	}
}

// Load reads a JSON config file and merges it over defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("config: decode: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Validate checks that the configuration values are sensible.
func (c *Config) Validate() error {
	if c.PortRange[0] < 1 || c.PortRange[1] > 65535 || c.PortRange[0] > c.PortRange[1] {
		return fmt.Errorf("config: invalid port_range [%d, %d]", c.PortRange[0], c.PortRange[1])
	}
	if c.Interval < 1 {
		return fmt.Errorf("config: interval_seconds must be >= 1")
	}
	if c.MaxHistory < 1 {
		return fmt.Errorf("config: max_history must be >= 1")
	}
	return nil
}
