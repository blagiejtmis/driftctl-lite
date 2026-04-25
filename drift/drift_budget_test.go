package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseBudgetReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc"},
		},
		Missing: []Resource{
			{ID: "sg-1", Type: "aws_security_group"},
			{ID: "sg-2", Type: "aws_security_group"},
		},
		Untracked: []Resource{
			{ID: "s3-1", Type: "aws_s3_bucket"},
		},
	}
}

func TestEvaluateBudget_NoDrift(t *testing.T) {
	report := Report{}
	cfg := DefaultBudgetConfig()
	br := EvaluateBudget(report, cfg)
	if BudgetHasViolations(br) {
		t.Error("expected no violations for empty report")
	}
}

func TestEvaluateBudget_WithinLimits(t *testing.T) {
	report := baseBudgetReport()
	cfg := DefaultBudgetConfig() // MaxMissing=5, MaxUntracked=10
	br := EvaluateBudget(report, cfg)
	if BudgetHasViolations(br) {
		t.Error("expected no violations when within limits")
	}
	if br.TotalMissing != 2 {
		t.Errorf("expected TotalMissing=2, got %d", br.TotalMissing)
	}
	if br.TotalUntracked != 1 {
		t.Errorf("expected TotalUntracked=1, got %d", br.TotalUntracked)
	}
}

func TestEvaluateBudget_MissingExceeded(t *testing.T) {
	report := baseBudgetReport()
	cfg := DefaultBudgetConfig()
	cfg.MaxMissing = 1
	br := EvaluateBudget(report, cfg)
	if !br.MissingExceeded {
		t.Error("expected MissingExceeded=true")
	}
	if !BudgetHasViolations(br) {
		t.Error("expected BudgetHasViolations=true")
	}
}

func TestEvaluateBudget_UntrackedExceeded(t *testing.T) {
	report := baseBudgetReport()
	cfg := DefaultBudgetConfig()
	cfg.MaxUntracked = 0
	br := EvaluateBudget(report, cfg)
	if !br.UntrackedExceeded {
		t.Error("expected UntrackedExceeded=true")
	}
}

func TestEvaluateBudget_PerTypeExceeded(t *testing.T) {
	report := baseBudgetReport()
	cfg := DefaultBudgetConfig()
	cfg.PerType = map[string]int{"aws_security_group": 1}
	br := EvaluateBudget(report, cfg)
	if len(br.Results) != 1 {
		t.Fatalf("expected 1 per-type result, got %d", len(br.Results))
	}
	if !br.Results[0].Exceeded {
		t.Error("expected per-type budget exceeded")
	}
	if !BudgetHasViolations(br) {
		t.Error("expected BudgetHasViolations=true")
	}
}

func TestFprintBudget_NoViolations(t *testing.T) {
	report := Report{}
	cfg := DefaultBudgetConfig()
	br := EvaluateBudget(report, cfg)
	var buf bytes.Buffer
	FprintBudget(&buf, br)
	out := buf.String()
	if !strings.Contains(out, "Drift Budget") {
		t.Error("expected header in output")
	}
	if strings.Contains(out, "EXCEEDED") {
		t.Error("did not expect EXCEEDED in output")
	}
}

func TestFprintBudget_WithViolations(t *testing.T) {
	report := baseBudgetReport()
	cfg := DefaultBudgetConfig()
	cfg.MaxMissing = 0
	cfg.PerType = map[string]int{"aws_security_group": 0}
	br := EvaluateBudget(report, cfg)
	var buf bytes.Buffer
	FprintBudget(&buf, br)
	out := buf.String()
	if !strings.Contains(out, "EXCEEDED") {
		t.Error("expected EXCEEDED in output")
	}
	if !strings.Contains(out, "aws_security_group") {
		t.Error("expected per-type row in output")
	}
}
