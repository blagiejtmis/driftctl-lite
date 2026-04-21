package drift

import (
	"bytes"
	"testing"
)

func baseComplianceReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{}},
			{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{}},
		},
	}
}

func TestEvaluateCompliance_NoFrameworks(t *testing.T) {
	results := EvaluateCompliance(baseComplianceReport(), nil)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEvaluateCompliance_AllPresent(t *testing.T) {
	fw := []ComplianceFramework{
		{Name: "baseline", Required: []string{"aws_vpc", "aws_security_group"}},
	}
	results := EvaluateCompliance(baseComplianceReport(), fw)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !results[0].Passing {
		t.Errorf("expected passing, got failing")
	}
	if len(results[0].Missing) != 0 {
		t.Errorf("expected no missing types")
	}
}

func TestEvaluateCompliance_MissingType(t *testing.T) {
	fw := []ComplianceFramework{
		{Name: "strict", Required: []string{"aws_vpc", "aws_iam_role"}},
	}
	results := EvaluateCompliance(baseComplianceReport(), fw)
	if results[0].Passing {
		t.Errorf("expected failing")
	}
	if len(results[0].Missing) != 1 || results[0].Missing[0] != "aws_iam_role" {
		t.Errorf("unexpected missing: %v", results[0].Missing)
	}
}

func TestComplianceHasFailures_True(t *testing.T) {
	results := []ComplianceResult{{Framework: "f", Passing: false}}
	if !ComplianceHasFailures(results) {
		t.Error("expected true")
	}
}

func TestComplianceHasFailures_False(t *testing.T) {
	results := []ComplianceResult{{Framework: "f", Passing: true}}
	if ComplianceHasFailures(results) {
		t.Error("expected false")
	}
}

func TestFprintCompliance_NoResults(t *testing.T) {
	var buf bytes.Buffer
	FprintCompliance(&buf, nil)
	if buf.Len() == 0 {
		t.Error("expected output")
	}
}

func TestFprintCompliance_WithResults(t *testing.T) {
	results := []ComplianceResult{
		{Framework: "cis", Passing: true, Missing: nil},
		{Framework: "pci", Passing: false, Missing: []string{"aws_kms_key"}},
	}
	var buf bytes.Buffer
	FprintCompliance(&buf, results)
	out := buf.String()
	if !contains(out, "PASS") {
		t.Error("expected PASS in output")
	}
	if !contains(out, "FAIL") {
		t.Error("expected FAIL in output")
	}
	if !contains(out, "aws_kms_key") {
		t.Error("expected missing type in output")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr([]string{sub}, sub) || func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
