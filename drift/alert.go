package drift

import (
	"fmt"
	"io"
	"time"
)

// AlertSeverity represents the urgency level of a drift alert.
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

// Alert represents a single drift alert generated from a report.
type Alert struct {
	Severity  AlertSeverity `json:"severity"`
	Type      string        `json:"type"`
	ResourceID string       `json:"resource_id"`
	Message   string        `json:"message"`
	Timestamp time.Time     `json:"timestamp"`
}

// AlertResult holds all alerts produced from a report.
type AlertResult struct {
	Alerts    []Alert `json:"alerts"`
	TotalInfo int     `json:"total_info"`
	TotalWarn int     `json:"total_warn"`
	TotalCrit int     `json:"total_crit"`
}

// EvaluateAlerts generates alerts from a drift report based on severity rules.
func EvaluateAlerts(report Report, critThreshold int) AlertResult {
	var result AlertResult
	now := time.Now().UTC()

	for _, r := range report.Missing {
		sev := SeverityWarning
		if len(report.Missing) >= critThreshold {
			sev = SeverityCritical
		}
		a := Alert{
			Severity:   sev,
			Type:       r.Type,
			ResourceID: r.ID,
			Message:    fmt.Sprintf("resource %s/%s is missing from live state", r.Type, r.ID),
			Timestamp:  now,
		}
		result.Alerts = append(result.Alerts, a)
		if sev == SeverityCritical {
			result.TotalCrit++
		} else {
			result.TotalWarn++
		}
	}

	for _, r := range report.Untracked {
		a := Alert{
			Severity:   SeverityInfo,
			Type:       r.Type,
			ResourceID: r.ID,
			Message:    fmt.Sprintf("resource %s/%s is untracked by IaC", r.Type, r.ID),
			Timestamp:  now,
		}
		result.Alerts = append(result.Alerts, a)
		result.TotalInfo++
	}

	return result
}

// AlertHasCritical returns true if any critical alerts exist.
func AlertHasCritical(result AlertResult) bool {
	return result.TotalCrit > 0
}

// FprintAlerts writes a human-readable alert summary to w.
func FprintAlerts(w io.Writer, result AlertResult) {
	if len(result.Alerts) == 0 {
		fmt.Fprintln(w, "No alerts.")
		return
	}
	fmt.Fprintf(w, "Alerts: %d critical, %d warning, %d info\n",
		result.TotalCrit, result.TotalWarn, result.TotalInfo)
	for _, a := range result.Alerts {
		fmt.Fprintf(w, "  [%s] %s: %s\n", a.Severity, a.ResourceID, a.Message)
	}
}
