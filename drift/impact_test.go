package drift

import (
	"bytes"
	"testing"
)

func baseImpactReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "sg-ok", Type: "aws_security_group"},
		},
		Missing: []Resource{
			{ID: "role-1", Type: "aws_iam_role"},
		},
		Untracked: []Resource{
			{ID: "bucket-1", Type: "aws_s3_bucket"},
		},
		Changed: []Resource{
			{ID: "alarm-1", Type: "aws_cloudwatch_alarm"},
		},
	}
}

func TestEvaluateImpact_NoResources(t *testing.T) {
	report := Report{}
	result := EvaluateImpact(report)
	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}

func TestEvaluateImpact_CorrectLevels(t *testing.T) {
	report := baseImpactReport()
	result := EvaluateImpact(report)

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}

	// IAM role should be critical
	if result.Results[0].Level != ImpactCritical {
		t.Errorf("expected first result to be critical, got %s", result.Results[0].Level)
	}
}

func TestEvaluateImpact_SortedByLevelDesc(t *testing.T) {
	report := baseImpactReport()
	result := EvaluateImpact(report)

	for i := 1; i < len(result.Results); i++ {
		prev := impactOrder(result.Results[i-1].Level)
		curr := impactOrder(result.Results[i].Level)
		if prev < curr {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}

func TestEvaluateImpact_UnknownTypeIsLow(t *testing.T) {
	report := Report{
		Missing: []Resource{
			{ID: "unknown-1", Type: "aws_some_new_service"},
		},
	}
	result := EvaluateImpact(report)
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if result.Results[0].Level != ImpactLow {
		t.Errorf("expected low impact for unknown type, got %s", result.Results[0].Level)
	}
}

func TestImpactHasCritical_True(t *testing.T) {
	r := ImpactReport{Results: []ImpactResult{{Level: ImpactCritical}}}
	if !ImpactHasCritical(r) {
		t.Error("expected true")
	}
}

func TestImpactHasCritical_False(t *testing.T) {
	r := ImpactReport{Results: []ImpactResult{{Level: ImpactLow}}}
	if ImpactHasCritical(r) {
		t.Error("expected false")
	}
}

func TestFprintImpact_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintImpact(&buf, ImpactReport{})
	if buf.Len() == 0 {
		t.Error("expected non-empty output for empty report")
	}
}

func TestFprintImpact_WithResults(t *testing.T) {
	var buf bytes.Buffer
	report := baseImpactReport()
	result := EvaluateImpact(report)
	FprintImpact(&buf, result)
	out := buf.String()
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("critical")) {
		t.Error("expected 'critical' in output")
	}
}
