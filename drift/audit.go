package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AuditEntry records a single scan event.
type AuditEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Command    string    `json:"command"`
	StateFile  string    `json:"state_file"`
	Total      int       `json:"total"`
	Missing    int       `json:"missing"`
	Untracked  int       `json:"untracked"`
	Changed    int       `json:"changed"`
	HasDrift   bool      `json:"has_drift"`
}

// AppendAudit appends an AuditEntry to the audit log file (newline-delimited JSON).
func AppendAudit(path string, entry AuditEntry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open %s: %w", path, err)
	}
	defer f.Close()
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// LoadAudit reads all AuditEntry records from a newline-delimited JSON file.
func LoadAudit(path string) ([]AuditEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read %s: %w", path, err)
	}
	var entries []AuditEntry
	for _, line := range splitLines(string(data)) {
		if line == "" {
			continue
		}
		var e AuditEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("audit: parse line: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// AuditEntryFromReport builds an AuditEntry from a Report and metadata.
func AuditEntryFromReport(cmd, stateFile string, r Report) AuditEntry {
	s := Summarize(r)
	return AuditEntry{
		Timestamp: time.Now().UTC(),
		Command:   cmd,
		StateFile: stateFile,
		Total:     s.Total,
		Missing:   s.Missing,
		Untracked: s.Untracked,
		Changed:   s.Changed,
		HasDrift:  s.HasDrift,
	}
}
