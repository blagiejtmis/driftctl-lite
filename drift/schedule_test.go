package drift

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func baseScheduleConfig() ScheduleConfig {
	return ScheduleConfig{
		Enabled:      true,
		IntervalMins: 30,
		StateFile:    "state.tfstate",
		OutputFormat: "json",
	}
}

func TestDefaultScheduleConfig(t *testing.T) {
	cfg := DefaultScheduleConfig()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if cfg.IntervalMins != 60 {
		t.Errorf("expected 60, got %d", cfg.IntervalMins)
	}
}

func TestSaveScheduleConfig_RoundTrip(t *testing.T) {
	cfg := baseScheduleConfig()
	path := filepath.Join(t.TempDir(), "schedule.json")
	if err := SaveScheduleConfig(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadScheduleConfig(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.IntervalMins != cfg.IntervalMins {
		t.Errorf("interval mismatch: got %d", loaded.IntervalMins)
	}
}

func TestLoadScheduleConfig_MissingFile(t *testing.T) {
	_, err := LoadScheduleConfig("/nonexistent/schedule.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadScheduleConfig_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "sched*.json")
	f.WriteString("{bad json")
	f.Close()
	_, err := LoadScheduleConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestIsDue_Disabled(t *testing.T) {
	cfg := baseScheduleConfig()
	cfg.Enabled = false
	if IsDue(cfg, time.Now()) {
		t.Error("disabled schedule should never be due")
	}
}

func TestIsDue_NeverRun(t *testing.T) {
	cfg := baseScheduleConfig()
	if !IsDue(cfg, time.Now()) {
		t.Error("should be due when never run")
	}
}

func TestIsDue_NotYet(t *testing.T) {
	cfg := baseScheduleConfig()
	cfg.LastRun = time.Now()
	if IsDue(cfg, time.Now()) {
		t.Error("should not be due immediately after last run")
	}
}

func TestIsDue_Overdue(t *testing.T) {
	cfg := baseScheduleConfig()
	cfg.LastRun = time.Now().Add(-2 * time.Hour)
	if !IsDue(cfg, time.Now()) {
		t.Error("should be due after interval elapsed")
	}
}

func TestFprintSchedule_Disabled(t *testing.T) {
	cfg := DefaultScheduleConfig()
	var buf bytes.Buffer
	FprintSchedule(&buf, cfg, time.Now())
	if !strings.Contains(buf.String(), "disabled") {
		t.Error("expected 'disabled' in output")
	}
}

func TestFprintSchedule_Due(t *testing.T) {
	cfg := baseScheduleConfig()
	cfg.LastRun = time.Now().Add(-2 * time.Hour)
	var buf bytes.Buffer
	FprintSchedule(&buf, cfg, time.Now())
	if !strings.Contains(buf.String(), "DUE") {
		t.Error("expected 'DUE' in output")
	}
}

func TestScheduleConfig_JSONFields(t *testing.T) {
	cfg := baseScheduleConfig()
	data, _ := json.Marshal(cfg)
	if !bytes.Contains(data, []byte("interval_mins")) {
		t.Error("expected interval_mins field in JSON")
	}
}
