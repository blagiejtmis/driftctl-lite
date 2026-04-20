package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ScheduleConfig defines when and how automated scans should run.
type ScheduleConfig struct {
	Enabled      bool          `json:"enabled"`
	IntervalMins int           `json:"interval_mins"`
	StateFile    string        `json:"state_file"`
	OutputFormat string        `json:"output_format"`
	LastRun      time.Time     `json:"last_run,omitempty"`
	NextRun      time.Time     `json:"next_run,omitempty"`
}

// DefaultScheduleConfig returns a ScheduleConfig with sensible defaults.
func DefaultScheduleConfig() ScheduleConfig {
	return ScheduleConfig{
		Enabled:      false,
		IntervalMins: 60,
		StateFile:    "terraform.tfstate",
		OutputFormat: "text",
	}
}

// LoadScheduleConfig reads a schedule config from a JSON file.
func LoadScheduleConfig(path string) (ScheduleConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultScheduleConfig(), fmt.Errorf("schedule config not found: %w", err)
	}
	var cfg ScheduleConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultScheduleConfig(), fmt.Errorf("invalid schedule config: %w", err)
	}
	return cfg, nil
}

// SaveScheduleConfig writes a ScheduleConfig to a JSON file.
func SaveScheduleConfig(path string, cfg ScheduleConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schedule config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// NextRunTime computes the next scheduled run time from a reference time.
func NextRunTime(from time.Time, intervalMins int) time.Time {
	return from.Add(time.Duration(intervalMins) * time.Minute)
}

// IsDue reports whether a scheduled scan is due relative to now.
func IsDue(cfg ScheduleConfig, now time.Time) bool {
	if !cfg.Enabled {
		return false
	}
	if cfg.LastRun.IsZero() {
		return true
	}
	return now.After(NextRunTime(cfg.LastRun, cfg.IntervalMins))
}
