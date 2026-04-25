package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseHealthReport() (Report, ScoreResult) {
	report := Report{
		Managed:   []Resource{{ID: "r1", Type: "aws_s3_bucket"}},
		Missing:   []Resource{},
		Untracked: []Resource{},
	}
	score := ScoreResult{Total: 1, Managed: 1, Percent: 100.0, Grade: "A"}
	return report, score
}

func TestEvaluateHealth_AllOK(t *testing.T) {
	report, score := baseHealthReport()
	h := EvaluateHealth(report, score)
	if h.Overall != HealthOK {
		t.Errorf("expected OK, got %s", h.Overall)
	}
	if len(h.Results) != 3 {
		t.Errorf("expected 3 checks, got %d", len(h.Results))
	}
}

func TestEvaluateHealth_MissingResourcesCritical(t *testing.T) {
	report, score := baseHealthReport()
	report.Missing = []Resource{{ID: "r2", Type: "aws_instance"}}
	h := EvaluateHealth(report, score)
	if h.Overall != HealthCritical {
		t.Errorf("expected CRITICAL, got %s", h.Overall)
	}
	if !HealthHasCritical(h) {
		t.Error("expected HealthHasCritical to return true")
	}
}

func TestEvaluateHealth_UntrackedWarning(t *testing.T) {
	report, score := baseHealthReport()
	report.Untracked = []Resource{{ID: "r3", Type: "aws_lambda_function"}}
	h := EvaluateHealth(report, score)
	if h.Overall != HealthWarning {
		t.Errorf("expected WARNING, got %s", h.Overall)
	}
	if HealthHasCritical(h) {
		t.Error("expected HealthHasCritical to return false")
	}
}

func TestEvaluateHealth_PoorGradeCritical(t *testing.T) {
	report, score := baseHealthReport()
	score.Grade = "F"
	score.Percent = 10.0
	h := EvaluateHealth(report, score)
	if h.Overall != HealthCritical {
		t.Errorf("expected CRITICAL for grade F, got %s", h.Overall)
	}
}

func TestEvaluateHealth_GradeBWarning(t *testing.T) {
	report, score := baseHealthReport()
	score.Grade = "B"
	score.Percent = 75.0
	h := EvaluateHealth(report, score)
	if h.Overall != HealthWarning {
		t.Errorf("expected WARNING for grade B, got %s", h.Overall)
	}
}

func TestFprintHealth_ContainsOverall(t *testing.T) {
	report, score := baseHealthReport()
	h := EvaluateHealth(report, score)
	var buf bytes.Buffer
	FprintHealth(&buf, h)
	out := buf.String()
	if !strings.Contains(out, "Health: OK") {
		t.Errorf("expected 'Health: OK' in output, got: %s", out)
	}
	if !strings.Contains(out, "drift_score") {
		t.Errorf("expected 'drift_score' check in output, got: %s", out)
	}
}

func TestHealthHasCritical_False(t *testing.T) {
	report, score := baseHealthReport()
	h := EvaluateHealth(report, score)
	if HealthHasCritical(h) {
		t.Error("expected no critical results")
	}
}
