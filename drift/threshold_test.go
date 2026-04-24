package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseThresholdScore(managed, total int) ScoreResult {
	pct := 0.0
	if total > 0 {
		pct = float64(managed) / float64(total) * 100.0
	}
	return ScoreResult{
		Total:      total,
		Managed:    managed,
		ManagedPct: pct,
	}
}

func TestEvaluateThreshold_Ok(t *testing.T) {
	cfg := DefaultThresholdConfig()
	score := baseThresholdScore(95, 100) // 5% drift
	r := EvaluateThreshold(score, cfg)
	if r.Status != "ok" {
		t.Errorf("expected ok, got %s", r.Status)
	}
	if r.DriftPct != 5.0 {
		t.Errorf("expected drift 5.0, got %.1f", r.DriftPct)
	}
}

func TestEvaluateThreshold_Warn(t *testing.T) {
	cfg := DefaultThresholdConfig() // warn=10, fail=25
	score := baseThresholdScore(85, 100) // 15% drift
	r := EvaluateThreshold(score, cfg)
	if r.Status != "warn" {
		t.Errorf("expected warn, got %s", r.Status)
	}
}

func TestEvaluateThreshold_Fail(t *testing.T) {
	cfg := DefaultThresholdConfig()
	score := baseThresholdScore(70, 100) // 30% drift
	r := EvaluateThreshold(score, cfg)
	if r.Status != "fail" {
		t.Errorf("expected fail, got %s", r.Status)
	}
}

func TestEvaluateThreshold_ZeroTotal(t *testing.T) {
	cfg := DefaultThresholdConfig()
	score := baseThresholdScore(0, 0)
	r := EvaluateThreshold(score, cfg)
	if r.Status != "ok" {
		t.Errorf("expected ok for empty, got %s", r.Status)
	}
}

func TestThresholdHasFailed_True(t *testing.T) {
	r := ThresholdResult{Status: "fail"}
	if !ThresholdHasFailed(r) {
		t.Error("expected true")
	}
}

func TestThresholdHasFailed_False(t *testing.T) {
	r := ThresholdResult{Status: "warn"}
	if ThresholdHasFailed(r) {
		t.Error("expected false")
	}
}

func TestFprintThreshold_Output(t *testing.T) {
	cfg := DefaultThresholdConfig()
	score := baseThresholdScore(80, 100)
	r := EvaluateThreshold(score, cfg)
	var buf bytes.Buffer
	FprintThreshold(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Threshold Check") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "20.0%") {
		t.Error("missing drift pct")
	}
	if !strings.Contains(out, r.Status) {
		t.Error("missing status")
	}
}
