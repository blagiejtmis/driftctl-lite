package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// RetentionConfig defines how long different data types should be kept.
type RetentionConfig struct {
	HistoryDays  int `json:"history_days"`
	AuditDays    int `json:"audit_days"`
	SnapshotDays int `json:"snapshot_days"`
	TrendDays    int `json:"trend_days"`
}

// RetentionResult summarises what was pruned.
type RetentionResult struct {
	Type    string
	Removed int
	Kept    int
}

// DefaultRetentionConfig returns sensible defaults.
func DefaultRetentionConfig() RetentionConfig {
	return RetentionConfig{
		HistoryDays:  90,
		AuditDays:    180,
		SnapshotDays: 30,
		TrendDays:    365,
	}
}

// LoadRetentionConfig reads a retention config from disk.
func LoadRetentionConfig(path string) (RetentionConfig, error) {
	cfg := DefaultRetentionConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// SaveRetentionConfig writes a retention config to disk.
func SaveRetentionConfig(path string, cfg RetentionConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ApplyRetention removes entries older than the configured thresholds.
// entries is a slice of timestamps; returns kept timestamps and removed count.
func ApplyRetention(entries []time.Time, maxDays int) ([]time.Time, int) {
	if maxDays <= 0 {
		return entries, 0
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -maxDays)
	var kept []time.Time
	removed := 0
	for _, t := range entries {
		if t.After(cutoff) {
			kept = append(kept, t)
		} else {
			removed++
		}
	}
	sort.Slice(kept, func(i, j int) bool { return kept[i].Before(kept[j]) })
	return kept, removed
}

// FprintRetention writes a human-readable retention report.
func FprintRetention(w io.Writer, results []RetentionResult) {
	fmt.Fprintln(w, "=== Retention Report ===")
	if len(results) == 0 {
		fmt.Fprintln(w, "  no data types evaluated")
		return
	}
	for _, r := range results {
		fmt.Fprintf(w, "  %-12s  kept: %d  removed: %d\n", r.Type, r.Kept, r.Removed)
	}
}

// RetentionHasRemovals returns true if any entries were pruned.
func RetentionHasRemovals(results []RetentionResult) bool {
	for _, r := range results {
		if r.Removed > 0 {
			return true
		}
	}
	return false
}
