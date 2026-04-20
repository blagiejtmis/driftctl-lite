package drift

import (
	"encoding/json"
	"fmt"
	"os"
)

// AlertConfig holds configuration for alert evaluation.
type AlertConfig struct {
	CriticalThreshold int    `json:"critical_threshold"`
	WebhookURL        string `json:"webhook_url,omitempty"`
	SlackChannel      string `json:"slack_channel,omitempty"`
}

// DefaultAlertConfig returns sensible defaults.
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		CriticalThreshold: 3,
	}
}

// LoadAlertConfig reads an AlertConfig from a JSON file.
func LoadAlertConfig(path string) (AlertConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultAlertConfig(), nil
		}
		return AlertConfig{}, fmt.Errorf("alert config: read: %w", err)
	}
	var cfg AlertConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AlertConfig{}, fmt.Errorf("alert config: parse: %w", err)
	}
	if cfg.CriticalThreshold <= 0 {
		cfg.CriticalThreshold = DefaultAlertConfig().CriticalThreshold
	}
	return cfg, nil
}

// SaveAlertConfig writes an AlertConfig to a JSON file.
func SaveAlertConfig(path string, cfg AlertConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("alert config: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("alert config: write: %w", err)
	}
	return nil
}

// AlertConfigFromEnv overrides config fields from environment variables.
func AlertConfigFromEnv(cfg AlertConfig) AlertConfig {
	if v := os.Getenv("DRIFTCTL_ALERT_WEBHOOK"); v != "" {
		cfg.WebhookURL = v
	}
	if v := os.Getenv("DRIFTCTL_ALERT_SLACK_CHANNEL"); v != "" {
		cfg.SlackChannel = v
	}
	return cfg
}
