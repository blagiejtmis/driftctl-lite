package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAlertConfig_MissingFile(t *testing.T) {
	cfg, err := LoadAlertConfig("/nonexistent/alert.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg.CriticalThreshold != 3 {
		t.Errorf("expected default threshold 3, got %d", cfg.CriticalThreshold)
	}
}

func TestLoadAlertConfig_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "alert.json")
	os.WriteFile(tmp, []byte("not-json"), 0644)
	_, err := LoadAlertConfig(tmp)
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestLoadAlertConfig_Valid(t *testing.T) {
	cfg := AlertConfig{CriticalThreshold: 7, WebhookURL: "https://example.com/hook"}
	data, _ := json.Marshal(cfg)
	tmp := filepath.Join(t.TempDir(), "alert.json")
	os.WriteFile(tmp, data, 0644)

	loaded, err := LoadAlertConfig(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.CriticalThreshold != 7 {
		t.Errorf("expected threshold 7, got %d", loaded.CriticalThreshold)
	}
	if loaded.WebhookURL != "https://example.com/hook" {
		t.Errorf("unexpected webhook url: %s", loaded.WebhookURL)
	}
}

func TestSaveAlertConfig_RoundTrip(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "alert.json")
	cfg := AlertConfig{CriticalThreshold: 5, SlackChannel: "#alerts"}
	if err := SaveAlertConfig(tmp, cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadAlertConfig(tmp)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.SlackChannel != "#alerts" {
		t.Errorf("expected #alerts, got %s", loaded.SlackChannel)
	}
}

func TestAlertConfigFromEnv(t *testing.T) {
	t.Setenv("DRIFTCTL_ALERT_WEBHOOK", "https://env-hook.example.com")
	t.Setenv("DRIFTCTL_ALERT_SLACK_CHANNEL", "#env-channel")
	cfg := AlertConfigFromEnv(DefaultAlertConfig())
	if cfg.WebhookURL != "https://env-hook.example.com" {
		t.Errorf("unexpected webhook: %s", cfg.WebhookURL)
	}
	if cfg.SlackChannel != "#env-channel" {
		t.Errorf("unexpected channel: %s", cfg.SlackChannel)
	}
}
