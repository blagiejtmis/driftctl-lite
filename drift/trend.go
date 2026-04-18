package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TrendEntry records a score snapshot at a point in time.
type TrendEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
	Grade     string    `json:"grade"`
	Managed   int       `json:"managed"`
	Missing   int       `json:"missing"`
	Untracked int       `json:"untracked"`
}

// TrendLog holds a series of trend entries.
type TrendLog struct {
	Entries []TrendEntry `json:"entries"`
}

// AppendTrend loads an existing trend log (if any), appends a new entry, and saves.
func AppendTrend(path string, entry TrendEntry) error {
	log, _ := LoadTrend(path)
	log.Entries = append(log.Entries, entry)
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("trend marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadTrend reads a trend log from disk.
func LoadTrend(path string) (TrendLog, error) {
	var log TrendLog
	data, err := os.ReadFile(path)
	if err != nil {
		return log, err
	}
	if err := json.Unmarshal(data, &log); err != nil {
		return log, fmt.Errorf("trend unmarshal: %w", err)
	}
	return log, nil
}

// TrendEntryFromScore builds a TrendEntry from a ScoreResult.
func TrendEntryFromScore(s ScoreResult) TrendEntry {
	return TrendEntry{
		Timestamp: time.Now().UTC(),
		Score:     s.Score,
		Grade:     s.Grade,
		Managed:   s.Managed,
		Missing:   s.Missing,
		Untracked: s.Untracked,
	}
}
