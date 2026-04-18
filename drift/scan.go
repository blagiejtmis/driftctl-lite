package drift

import (
	"encoding/json"
	"fmt"
	"os"\n)

// Resource represents a cloud resource parsed from state.
type Resource struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
}

// DriftResult holds the outcome of a drift scan.
type DriftResult struct {
	Managed []Resource `json:"managed"`
	Drifted []Resource `json:"drifted"`
	Missing []Resource `json:"missing"`
}

// Scan reads the given Terraform state file and detects drift.
func Scan(statePath string) (*DriftResult, error) {
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", statePath, err)
	}

	var state tfState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}

	resources := parseResources(state)
	return detectDrift(resources)
}
