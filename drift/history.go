package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// HistoryEntry records a scan result at a point in time.
type HistoryEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	ManagedCount   int   `json:"managed_count"`
	MissingCount   int   `json:"missing_count"`
	UntrackedCount int   `json:"untracked_count"`
	DriftScore     float64 `json:"drift_score"`
	Label          string  `json:"label,omitempty"`
}

// AppendHistory appends an entry to a JSON-lines history file.
func AppendHistory(path string, entry HistoryEntry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open history file: %w", err)
	}
	defer f.Close()
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// LoadHistory reads all history entries from a JSON-lines file.
func LoadHistory(path string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read history file: %w", err)
	}
	var entries []HistoryEntry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e HistoryEntry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("parse history entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// splitLines splits a byte slice by newline.
func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

// HistoryEntryFromScore builds a HistoryEntry from a ScoreResult.
func HistoryEntryFromScore(s ScoreResult, label string) HistoryEntry {
	return HistoryEntry{
		Timestamp:      time.Now().UTC(),
		ManagedCount:   s.Managed,
		MissingCount:   s.Missing,
		UntrackedCount: s.Untracked,
		DriftScore:     s.Score,
		Label:          label,
	}
}
