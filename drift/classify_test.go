package drift

import (
	"bytes"
	"testing"
)

func baseClassifyReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "sg-1", Type: "aws_security_group"},
		},
		Missing: []Resource{
			{ID: "role-1", Type: "aws_iam_role"},
			{ID: "bucket-1", Type: "aws_s3_bucket"},
		},
		Untracked: []Resource{
			{ID: "ec2-1", Type: "aws_instance"},
		},
	}
}

func TestClassify_NoResources(t *testing.T) {
	r := Classify(Report{})
	if len(r.Results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(r.Results))
	}
}

func TestClassify_CorrectSeverities(t *testing.T) {
	r := Classify(baseClassifyReport())
	if len(r.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(r.Results))
	}
	// first result should be critical (IAM role)
	if r.Results[0].Severity != ClassCritical {
		t.Errorf("expected critical, got %s", r.Results[0].Severity)
	}
}

func TestClassify_SortedDescending(t *testing.T) {
	r := Classify(baseClassifyReport())
	for i := 1; i < len(r.Results); i++ {
		if severityOrder[r.Results[i].Severity] > severityOrder[r.Results[i-1].Severity] {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}

func TestClassifyHasCritical_True(t *testing.T) {
	r := Classify(baseClassifyReport())
	if !ClassifyHasCritical(r) {
		t.Error("expected ClassifyHasCritical to return true")
	}
}

func TestClassifyHasCritical_False(t *testing.T) {
	report := Report{
		Untracked: []Resource{{ID: "ec2-1", Type: "aws_instance"}},
	}
	r := Classify(report)
	if ClassifyHasCritical(r) {
		t.Error("expected ClassifyHasCritical to return false")
	}
}

func TestFprintClassify_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintClassify(&buf, ClassifyReport{})
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestFprintClassify_WithResults(t *testing.T) {
	r := Classify(baseClassifyReport())
	var buf bytes.Buffer
	FprintClassify(&buf, r)
	out := buf.String()
	if len(out) == 0 {
		t.Error("expected output")
	}
	for _, res := range r.Results {
		if !bytes.Contains(buf.Bytes(), []byte(res.Resource.ID)) {
			t.Errorf("expected output to contain resource ID %s", res.Resource.ID)
		}
	}
}
