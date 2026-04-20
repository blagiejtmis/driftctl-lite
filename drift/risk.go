package drift

import (
	"fmt"
	"io"
	"sort"
)

// RiskLevel represents the severity of a drift risk.
type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// RiskEntry describes the assessed risk for a single resource.
type RiskEntry struct {
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Reason       string    `json:"reason"`
	Level        RiskLevel `json:"level"`
}

// RiskReport holds all risk entries produced from a drift report.
type RiskReport struct {
	Entries []RiskEntry `json:"entries"`
}

// riskWeights assigns a numeric weight to each risk level for sorting.
var riskWeights = map[RiskLevel]int{
	RiskCritical: 4,
	RiskHigh:     3,
	RiskMedium:   2,
	RiskLow:      1,
}

// highRiskTypes lists resource types that carry elevated risk when drifted.
var highRiskTypes = map[string]bool{
	"aws_iam_role":          true,
	"aws_iam_policy":        true,
	"aws_security_group":    true,
	"aws_s3_bucket":         true,
	"aws_kms_key":           true,
	"aws_secretsmanager_secret": true,
}

// EvaluateRisk inspects a Report and assigns a risk level to each drifted
// or untracked resource based on its type and drift status.
func EvaluateRisk(report Report) RiskReport {
	var entries []RiskEntry

	for _, r := range report.Missing {
		level := RiskMedium
		reason := "resource defined in IaC is missing from live infrastructure"
		if highRiskTypes[r.Type] {
			level = RiskCritical
			reason = fmt.Sprintf("high-risk resource type %q is missing from live infrastructure", r.Type)
		}
		entries = append(entries, RiskEntry{
			ResourceID:   r.ID,
			ResourceType: r.Type,
			Reason:       reason,
			Level:        level,
		})
	}

	for _, r := range report.Untracked {
		level := RiskLow
		reason := "resource exists in live infrastructure but is not tracked by IaC"
		if highRiskTypes[r.Type] {
			level = RiskHigh
			reason = fmt.Sprintf("high-risk resource type %q is untracked by IaC", r.Type)
		}
		entries = append(entries, RiskEntry{
			ResourceID:   r.ID,
			ResourceType: r.Type,
			Reason:       reason,
			Level:        level,
		})
	}

	// Sort by descending risk weight, then by resource ID for determinism.
	sort.Slice(entries, func(i, j int) bool {
		wi := riskWeights[entries[i].Level]
		wj := riskWeights[entries[j].Level]
		if wi != wj {
			return wi > wj
		}
		return entries[i].ResourceID < entries[j].ResourceID
	})

	return RiskReport{Entries: entries}
}

// RiskHasCritical returns true when any entry in the report is critical.
func RiskHasCritical(r RiskReport) bool {
	for _, e := range r.Entries {
		if e.Level == RiskCritical {
			return true
		}
	}
	return false
}

// FprintRisk writes a human-readable risk report to w.
func FprintRisk(w io.Writer, r RiskReport) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "Risk assessment: no drift detected, all resources are low risk.")
		return
	}
	fmt.Fprintf(w, "Risk assessment (%d issue(s)):\n", len(r.Entries))
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  [%-8s] %s (%s)\n    %s\n", e.Level, e.ResourceID, e.ResourceType, e.Reason)
	}
}
