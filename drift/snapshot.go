package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures a point-in-time view of a scan report.
type Snapshot struct {
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label"`
	Report    Report    `json:"report"`
}

// SaveSnapshot writes a snapshot to the given file path as JSON.
func SaveSnapshot(path, label string, report Report) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Report:    report,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot write: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from the given file path.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot parse: %w", err)
	}
	return snap, nil
}
