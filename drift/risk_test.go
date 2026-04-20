package drift

import (
	"bytes"
	"testing"
)

func baseRiskResources() []Resource {
	return []Resource{
		{ID: "r1", Type: "aws_s3_bucket", Attributes: map[string]interface{}{"public": true}},
		{ID: "r2", Type: "aws_iam_role", Attributes: map[string]interface{}{"policy": "AdministratorAccess"}},
		{ID: "r3", Type: "aws_instance", Attributes: map[string]interface{}{"instance_type": "t2.micro"}},
	}
}

func TestEvaluateRisk_NoResources(t *testing.T) {
	results := EvaluateRisk([]Resource{})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEvaluateRisk_KnownHighRisk(t *testing.T) {
	res := baseRiskResources()
	results := EvaluateRisk(res)
	found := false
	for _, r := range results {
		if r.ResourceID == "r2" && r.Level == "critical" {
			found = true
		}
	}
	if !found {
		t.Error("expected critical risk for aws_iam_role with AdministratorAccess")
	}
}

func TestRiskHasCritical_True(t *testing.T) {
	results := EvaluateRisk(baseRiskResources())
	if !RiskHasCritical(results) {
		t.Error("expected RiskHasCritical to be true")
	}
}

func TestRiskHasCritical_False(t *testing.T) {
	results := []RiskResult{}
	if RiskHasCritical(results) {
		t.Error("expected RiskHasCritical to be false")
	}
}

func TestFprintRisk_NoResults(t *testing.T) {
	var buf bytes.Buffer
	FprintRisk(&buf, []RiskResult{})
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output even with no results")
	}
}

func TestFprintRisk_WithResults(t *testing.T) {
	results := EvaluateRisk(baseRiskResources())
	var buf bytes.Buffer
	FprintRisk(&buf, results)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output")
	}
}
