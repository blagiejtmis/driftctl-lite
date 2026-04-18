package drift

import (
	"encoding/json"
	"fmt"
	"os"
)

const defaultNotifyConfigFile = ".driftctl-notify.json"

// LoadNotifyConfig reads a NotifyConfig from a JSON file.
func LoadNotifyConfig(path string) (NotifyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NotifyConfig{}, fmt.Errorf("notify config file not found: %s", path)
		}
		return NotifyConfig{}, fmt.Errorf("reading notify config: %w", err)
	}
	var cfg NotifyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return NotifyConfig{}, fmt.Errorf("parsing notify config: %w", err)
	}
	return cfg, nil
}

// SaveNotifyConfig writes a NotifyConfig to a JSON file.
func SaveNotifyConfig(path string, cfg NotifyConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling notify config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing notify config: %w", err)
	}
	return nil
}

// NotifyConfigFromEnv builds a NotifyConfig from environment variables.
func NotifyConfigFromEnv() NotifyConfig {
	return NotifyConfig{
		WebhookURL: os.Getenv("DRIFTCTL_WEBHOOK_URL"),
		SlackToken: os.Getenv("DRIFTCTL_SLACK_TOKEN"),
		Channel:    os.Getenv("DRIFTCTL_SLACK_CHANNEL"),
		OnDriftOnly: os.Getenv("DRIFTCTL_ON_DRIFT_ONLY") != "false",
	}
}
