package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved drift report used as a reference point.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Resources []Resource `json:"resources"`
}

// SaveBaseline writes the managed resources from a report to a baseline file.
func SaveBaseline(path string, report Report) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Resources: report.Managed,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// LoadBaseline reads a baseline file and returns the stored resources.
func LoadBaseline(path string) (Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return Baseline{}, fmt.Errorf("baseline: parse %s: %w", path, err)
	}
	return b, nil
}
