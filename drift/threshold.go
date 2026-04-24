package drift

import (
	"fmt"
	"io"
)

// ThresholdConfig defines drift percentage thresholds for pass/warn/fail.
type ThresholdConfig struct {
	WarnPct float64 `json:"warn_pct"`
	FailPct float64 `json:"fail_pct"`
}

// ThresholdResult holds the evaluation outcome for a single threshold check.
type ThresholdResult struct {
	DriftPct float64
	WarnPct  float64
	FailPct  float64
	Status   string // "ok", "warn", "fail"
}

// DefaultThresholdConfig returns sensible defaults.
func DefaultThresholdConfig() ThresholdConfig {
	return ThresholdConfig{
		WarnPct: 10.0,
		FailPct: 25.0,
	}
}

// EvaluateThreshold checks the drift percentage against warn/fail thresholds.
func EvaluateThreshold(score ScoreResult, cfg ThresholdConfig) ThresholdResult {
	driftPct := 100.0 - score.ManagedPct
	status := "ok"
	if driftPct >= cfg.FailPct {
		status = "fail"
	} else if driftPct >= cfg.WarnPct {
		status = "warn"
	}
	return ThresholdResult{
		DriftPct: driftPct,
		WarnPct:  cfg.WarnPct,
		FailPct:  cfg.FailPct,
		Status:   status,
	}
}

// ThresholdHasFailed returns true when the result status is "fail".
func ThresholdHasFailed(r ThresholdResult) bool {
	return r.Status == "fail"
}

// FprintThreshold writes a human-readable threshold evaluation to w.
func FprintThreshold(w io.Writer, r ThresholdResult) {
	icon := "✔"
	switch r.Status {
	case "warn":
		icon = "⚠"
	case "fail":
		icon = "✘"
	}
	fmt.Fprintf(w, "Threshold Check: %s\n", icon)
	fmt.Fprintf(w, "  Drift:   %.1f%%\n", r.DriftPct)
	fmt.Fprintf(w, "  Warn at: %.1f%%\n", r.WarnPct)
	fmt.Fprintf(w, "  Fail at: %.1f%%\n", r.FailPct)
	fmt.Fprintf(w, "  Status:  %s\n", r.Status)
}
