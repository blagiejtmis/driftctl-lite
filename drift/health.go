package drift

import (
	"fmt"
	"io"
	"sort"
)

// HealthStatus represents the overall health level of a drift scan.
type HealthStatus string

const (
	HealthOK      HealthStatus = "OK"
	HealthWarning HealthStatus = "WARNING"
	HealthCritical HealthStatus = "CRITICAL"
)

// HealthResult holds the evaluated health for a single check.
type HealthResult struct {
	Check   string
	Status  HealthStatus
	Message string
}

// HealthReport aggregates all health check results.
type HealthReport struct {
	Results []HealthResult
	Overall HealthStatus
}

// EvaluateHealth runs a set of health checks against a drift Report and Score.
func EvaluateHealth(report Report, score ScoreResult) HealthReport {
	var results []HealthResult

	// Check 1: untracked resources
	if len(report.Untracked) > 0 {
		results = append(results, HealthResult{
			Check:   "untracked_resources",
			Status:  HealthWarning,
			Message: fmt.Sprintf("%d untracked resource(s) found", len(report.Untracked)),
		})
	} else {
		results = append(results, HealthResult{
			Check:   "untracked_resources",
			Status:  HealthOK,
			Message: "no untracked resources",
		})
	}

	// Check 2: missing resources
	if len(report.Missing) > 0 {
		results = append(results, HealthResult{
			Check:   "missing_resources",
			Status:  HealthCritical,
			Message: fmt.Sprintf("%d missing resource(s) detected", len(report.Missing)),
		})
	} else {
		results = append(results, HealthResult{
			Check:   "missing_resources",
			Status:  HealthOK,
			Message: "no missing resources",
		})
	}

	// Check 3: drift score
	switch {
	case score.Grade == "A":
		results = append(results, HealthResult{
			Check:   "drift_score",
			Status:  HealthOK,
			Message: fmt.Sprintf("drift score %.1f%% (grade %s)", score.Percent, score.Grade),
		})
	case score.Grade == "B" || score.Grade == "C":
		results = append(results, HealthResult{
			Check:   "drift_score",
			Status:  HealthWarning,
			Message: fmt.Sprintf("drift score %.1f%% (grade %s)", score.Percent, score.Grade),
		})
	default:
		results = append(results, HealthResult{
			Check:   "drift_score",
			Status:  HealthCritical,
			Message: fmt.Sprintf("drift score %.1f%% (grade %s)", score.Percent, score.Grade),
		})
	}

	overall := overallHealth(results)
	return HealthReport{Results: results, Overall: overall}
}

func overallHealth(results []HealthResult) HealthStatus {
	statuses := map[HealthStatus]int{}
	for _, r := range results {
		statuses[r.Status]++
	}
	if statuses[HealthCritical] > 0 {
		return HealthCritical
	}
	if statuses[HealthWarning] > 0 {
		return HealthWarning
	}
	return HealthOK
}

// HealthHasCritical returns true if any check is CRITICAL.
func HealthHasCritical(h HealthReport) bool {
	for _, r := range h.Results {
		if r.Status == HealthCritical {
			return true
		}
	}
	return false
}

// FprintHealth writes a human-readable health summary to w.
func FprintHealth(w io.Writer, h HealthReport) {
	fmt.Fprintf(w, "Health: %s\n", h.Overall)

	// sort for deterministic output
	sorted := make([]HealthResult, len(h.Results))
	copy(sorted, h.Results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Check < sorted[j].Check
	})

	for _, r := range sorted {
		fmt.Fprintf(w, "  [%s] %s: %s\n", r.Status, r.Check, r.Message)
	}
}
