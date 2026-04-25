package drift

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baseRetentionConfig() RetentionConfig {
	return RetentionConfig{
		HistoryDays:  90,
		AuditDays:    180,
		SnapshotDays: 30,
		TrendDays:    365,
	}
}

func TestDefaultRetentionConfig(t *testing.T) {
	cfg := DefaultRetentionConfig()
	if cfg.HistoryDays != 90 {
		t.Errorf("expected HistoryDays=90, got %d", cfg.HistoryDays)
	}
	if cfg.AuditDays != 180 {
		t.Errorf("expected AuditDays=180, got %d", cfg.AuditDays)
	}
}

func TestSaveRetentionConfig_RoundTrip(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "retention.json")
	cfg := baseRetentionConfig()
	if err := SaveRetentionConfig(tmp, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRetentionConfig(tmp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.HistoryDays != cfg.HistoryDays || loaded.TrendDays != cfg.TrendDays {
		t.Errorf("round-trip mismatch: %+v", loaded)
	}
}

func TestLoadRetentionConfig_MissingFile(t *testing.T) {
	cfg, err := LoadRetentionConfig("/nonexistent/retention.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg.HistoryDays != 90 {
		t.Errorf("expected defaults, got %+v", cfg)
	}
}

func TestLoadRetentionConfig_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not-json"), 0o644)
	_, err := LoadRetentionConfig(tmp)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApplyRetention_RemovesOld(t *testing.T) {
	now := time.Now().UTC()
	entries := []time.Time{
		now.AddDate(0, 0, -100),
		now.AddDate(0, 0, -50),
		now.AddDate(0, 0, -5),
	}
	kept, removed := ApplyRetention(entries, 60)
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
	if len(kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(kept))
	}
}

func TestApplyRetention_ZeroMaxDays(t *testing.T) {
	now := time.Now().UTC()
	entries := []time.Time{now.AddDate(0, 0, -500)}
	kept, removed := ApplyRetention(entries, 0)
	if removed != 0 || len(kept) != 1 {
		t.Errorf("zero maxDays should keep all: kept=%d removed=%d", len(kept), removed)
	}
}

func TestFprintRetention_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintRetention(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("no data types")) {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestFprintRetention_WithResults(t *testing.T) {
	results := []RetentionResult{
		{Type: "history", Removed: 3, Kept: 10},
		{Type: "audit", Removed: 0, Kept: 5},
	}
	var buf bytes.Buffer
	FprintRetention(&buf, results)
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("history")) {
		t.Errorf("expected 'history' in output, got %q", out)
	}
	if !bytes.Contains([]byte(out), []byte("removed: 3")) {
		t.Errorf("expected 'removed: 3' in output, got %q", out)
	}
}

func TestRetentionHasRemovals(t *testing.T) {
	if RetentionHasRemovals([]RetentionResult{{Type: "x", Removed: 0}}) {
		t.Error("expected false")
	}
	if !RetentionHasRemovals([]RetentionResult{{Type: "x", Removed: 1}}) {
		t.Error("expected true")
	}
}

func TestRetentionConfig_JSON(t *testing.T) {
	cfg := baseRetentionConfig()
	data, _ := json.Marshal(cfg)
	var out RetentionConfig
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.SnapshotDays != cfg.SnapshotDays {
		t.Errorf("SnapshotDays mismatch")
	}
}
