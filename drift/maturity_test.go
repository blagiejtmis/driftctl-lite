package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseMaturityReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "a", Type: "aws_s3_bucket"},
			{ID: "b", Type: "aws_s3_bucket"},
			{ID: "c", Type: "aws_iam_role"},
		},
		Missing: []Resource{
			{ID: "d", Type: "aws_s3_bucket"},
		},
		Untracked: []Resource{
			{ID: "e", Type: "aws_lambda_function"},
		},
	}
}

func TestEvaluateMaturity_Empty(t *testing.T) {
	r := EvaluateMaturity(Report{})
	if len(r.Results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(r.Results))
	}
	if r.Overall != MaturityInitial {
		t.Errorf("expected Initial for empty report, got %s", r.Overall)
	}
}

func TestEvaluateMaturity_CorrectCoverage(t *testing.T) {
	r := EvaluateMaturity(baseMaturityReport())
	var s3 *MaturityResult
	for i := range r.Results {
		if r.Results[i].Type == "aws_s3_bucket" {
			s3 = &r.Results[i]
		}
	}
	if s3 == nil {
		t.Fatal("aws_s3_bucket not found in results")
	}
	// 2 managed out of 3 total => 66.7%
	if s3.Total != 3 {
		t.Errorf("expected total 3, got %d", s3.Total)
	}
	if s3.Managed != 2 {
		t.Errorf("expected managed 2, got %d", s3.Managed)
	}
	if s3.Level != MaturityDefined {
		t.Errorf("expected Defined, got %s", s3.Level)
	}
}

func TestEvaluateMaturity_UntrackedCountsAsUnmanaged(t *testing.T) {
	r := EvaluateMaturity(baseMaturityReport())
	var lambda *MaturityResult
	for i := range r.Results {
		if r.Results[i].Type == "aws_lambda_function" {
			lambda = &r.Results[i]
		}
	}
	if lambda == nil {
		t.Fatal("aws_lambda_function not found")
	}
	if lambda.Coverage != 0 {
		t.Errorf("expected 0%% coverage, got %.1f", lambda.Coverage)
	}
	if lambda.Level != MaturityInitial {
		t.Errorf("expected Initial, got %s", lambda.Level)
	}
}

func TestMaturityHasCritical_True(t *testing.T) {
	r := MaturityReport{Overall: MaturityInitial}
	if !MaturityHasCritical(r) {
		t.Error("expected critical for Initial")
	}
}

func TestMaturityHasCritical_False(t *testing.T) {
	r := MaturityReport{Overall: MaturityOptimizing}
	if MaturityHasCritical(r) {
		t.Error("expected no critical for Optimizing")
	}
}

func TestFprintMaturity_ContainsTypes(t *testing.T) {
	r := EvaluateMaturity(baseMaturityReport())
	var buf bytes.Buffer
	FprintMaturity(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Error("expected aws_s3_bucket in output")
	}
	if !strings.Contains(out, "Overall") {
		t.Error("expected Overall in output")
	}
}

func TestFprintMaturity_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintMaturity(&buf, MaturityReport{Overall: MaturityInitial})
	if !strings.Contains(buf.String(), "No resource types") {
		t.Error("expected empty message")
	}
}
