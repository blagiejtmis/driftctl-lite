package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseDriftMapReport() Report {
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

func TestBuildDriftMap_Empty(t *testing.T) {
	dm := BuildDriftMap(Report{})
	if len(dm) != 0 {
		t.Fatalf("expected empty drift map, got %d entries", len(dm))
	}
}

func TestBuildDriftMap_CorrectCounts(t *testing.T) {
	dm := BuildDriftMap(baseDriftMapReport())

	index := map[string]DriftMapEntry{}
	for _, e := range dm {
		index[e.Type] = e
	}

	s3 := index["aws_s3_bucket"]
	if s3.Total != 3 || s3.Missing != 1 {
		t.Errorf("aws_s3_bucket: got total=%d missing=%d", s3.Total, s3.Missing)
	}

	inst := index["aws_instance"]
	if inst.Total != 2 || inst.Untracked != 1 {
		t.Errorf("aws_instance: got total=%d untracked=%d", inst.Total, inst.Untracked)
	}

	lambda := index["aws_lambda_function"]
	if lambda.Total != 1 || lambda.Changed != 1 {
		t.Errorf("aws_lambda_function: got total=%d changed=%d", lambda.Total, lambda.Changed)
	}
}

func TestBuildDriftMap_SortedByDriftPctDesc(t *testing.T) {
	dm := BuildDriftMap(baseDriftMapReport())
	for i := 1; i < len(dm); i++ {
		if dm[i].DriftPct > dm[i-1].DriftPct {
			t.Errorf("drift map not sorted descending at index %d", i)
		}
	}
}

func TestDriftMapHasDrift_True(t *testing.T) {
	dm := BuildDriftMap(baseDriftMapReport())
	if !DriftMapHasDrift(dm) {
		t.Error("expected drift to be detected")
	}
}

func TestDriftMapHasDrift_False(t *testing.T) {
	r := Report{
		Managed: []Resource{{ID: "x", Type: "aws_s3_bucket"}},
	}
	dm := BuildDriftMap(r)
	if DriftMapHasDrift(dm) {
		t.Error("expected no drift")
	}
}

func TestFprintDriftMap_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintDriftMap(&buf, DriftMap{})
	if !strings.Contains(buf.String(), "no resources") {
		t.Errorf("expected 'no resources' message, got: %s", buf.String())
	}
}

func TestFprintDriftMap_OutputContainsTypes(t *testing.T) {
	dm := BuildDriftMap(baseDriftMapReport())
	var buf bytes.Buffer
	FprintDriftMap(&buf, dm)
	out := buf.String()
	for _, typ := range []string{"aws_s3_bucket", "aws_instance", "aws_lambda_function"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected output to contain %q", typ)
		}
	}
}
