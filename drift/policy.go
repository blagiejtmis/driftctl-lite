package drift

import (
	"encoding/json"
	"fmt"
	"os"
)

// PolicyRule defines a rule that can fail or warn on drift conditions.
type PolicyRule struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Severity string `json:"severity"` // "error" or "warn"
	Message  string `json:"message"`
}

// PolicyFile holds a list of rules loaded from a JSON policy file.
type PolicyFile struct {
	Rules []PolicyRule `json:"rules"`
}

// PolicyResult holds the outcome of evaluating a single rule.
type PolicyResult struct {
	Rule      PolicyRule
	Violations []Resource
}

// LoadPolicy reads a policy file from disk.
func LoadPolicy(path string) (*PolicyFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read %s: %w", path, err)
	}
	var pf PolicyFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("policy: parse %s: %w", path, err)
	}
	return &pf, nil
}

// EvaluatePolicy checks a Report against the loaded policy rules.
func EvaluatePolicy(pf *PolicyFile, report Report) []PolicyResult {
	var results []PolicyResult
	for _, rule := range pf.Rules {
		var violations []Resource
		for _, r := range report.Missing {
			if rule.Type == "" || rule.Type == r.Type {
				violations = append(violations, r)
			}
		}
		if len(violations) > 0 {
			results = append(results, PolicyResult{Rule: rule, Violations: violations})
		}
	}
	return results
}
