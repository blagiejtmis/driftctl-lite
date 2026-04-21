package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseRollupReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "a", Type: "aws_s3_bucket"},
			{ID: "b", Type: "aws_s3_bucket"},
			{ID: "c", Type: "aws_instance"},
		},
		Missing: []Resource{
			{ID: "d", Type: "aws_s3_bucket"},
		},
		Untracked: []Resource{
			{ID: "e", Type: "aws_instance"},
		},
		Changed: []Resource{
			{ID: "f", Type: "aws_lambda_function"},
		},
	}
}

func TestRollup_Empty(t *testing.T) {
	entries := Rollup(Report{})
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestRollup_CorrectCounts(t *testing.T) {
	entries := Rollup(baseRollupReport())
	byType := map[string]RollupEntry{}
	for _, e := range entries {
		byType[e.Type] = e
	}

	s3 := byType["aws_s3_bucket"]
	if s3.Managed != 2 || s3.Missing != 1 || s3.Total != 3 {
		t.Errorf("unexpected s3 counts: %+v", s3)
	}

	ec2 := byType["aws_instance"]
	if ec2.Managed != 1 || ec2.Untracked != 1 || ec2.Total != 2 {
		t.Errorf("unexpected ec2 counts: %+v", ec2)
	}

	lambda := byType["aws_lambda_function"]
	if lambda.Changed != 1 || lambda.Total != 1 {
		t.Errorf("unexpected lambda counts: %+v", lambda)
	}
}

func TestRollup_SortedByType(t *testing.T) {
	entries := Rollup(baseRollupReport())
	for i := 1; i < len(entries); i++ {
		if entries[i].Type < entries[i-1].Type {
			t.Errorf("entries not sorted: %s before %s", entries[i-1].Type, entries[i].Type)
		}
	}
}

func TestFprintRollup_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintRollup(&buf, []RollupEntry{})
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintRollup_ContainsTypes(t *testing.T) {
	entries := Rollup(baseRollupReport())
	var buf bytes.Buffer
	FprintRollup(&buf, entries)
	out := buf.String()
	for _, typ := range []string{"aws_s3_bucket", "aws_instance", "aws_lambda_function"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected %q in output", typ)
		}
	}
}

func TestRollupHasDrift_True(t *testing.T) {
	entries := Rollup(baseRollupReport())
	if !RollupHasDrift(entries) {
		t.Error("expected drift to be detected")
	}
}

func TestRollupHasDrift_False(t *testing.T) {
	r := Report{
		Managed: []Resource{{ID: "x", Type: "aws_s3_bucket"}},
	}
	if RollupHasDrift(Rollup(r)) {
		t.Error("expected no drift")
	}
}
