package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseDensityReport() Report {
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
			{ID: "f", Type: "aws_lambda_function"},
		},
	}
}

func TestBuildDriftDensity_Empty(t *testing.T) {
	r := BuildDriftDensity(Report{})
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.OverallDensity != 0 {
		t.Fatalf("expected 0 overall density, got %f", r.OverallDensity)
	}
}

func TestBuildDriftDensity_CorrectCounts(t *testing.T) {
	r := BuildDriftDensity(baseDensityReport())
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
	byType := map[string]DriftDensityEntry{}
	for _, e := range r.Entries {
		byType[e.Type] = e
	}
	s3 := byType["aws_s3_bucket"]
	if s3.Total != 3 || s3.Drifted != 1 {
		t.Errorf("aws_s3_bucket: want total=3 drifted=1, got total=%d drifted=%d", s3.Total, s3.Drifted)
	}
	inst := byType["aws_instance"]
	if inst.Total != 2 || inst.Drifted != 1 {
		t.Errorf("aws_instance: want total=2 drifted=1, got total=%d drifted=%d", inst.Total, inst.Drifted)
	}
	lambda := byType["aws_lambda_function"]
	if lambda.Total != 1 || lambda.Drifted != 1 {
		t.Errorf("aws_lambda_function: want total=1 drifted=1, got total=%d drifted=%d", lambda.Total, lambda.Drifted)
	}
}

func TestBuildDriftDensity_SortedByDensityDesc(t *testing.T) {
	r := BuildDriftDensity(baseDensityReport())
	for i := 1; i < len(r.Entries); i++ {
		if r.Entries[i].DensityPct > r.Entries[i-1].DensityPct {
			t.Errorf("entries not sorted desc at index %d", i)
		}
	}
}

func TestBuildDriftDensity_OverallDensity(t *testing.T) {
	r := BuildDriftDensity(baseDensityReport())
	// 6 total resources, 3 drifted => 50%
	if r.OverallDensity < 49.9 || r.OverallDensity > 50.1 {
		t.Errorf("expected ~50%% overall density, got %.2f", r.OverallDensity)
	}
}

func TestDriftDensityHasDrift_True(t *testing.T) {
	r := BuildDriftDensity(baseDensityReport())
	if !DriftDensityHasDrift(r) {
		t.Error("expected HasDrift=true")
	}
}

func TestDriftDensityHasDrift_False(t *testing.T) {
	r := BuildDriftDensity(Report{
		Managed: []Resource{{ID: "x", Type: "aws_s3_bucket"}},
	})
	if DriftDensityHasDrift(r) {
		t.Error("expected HasDrift=false")
	}
}

func TestFprintDensity_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintDensity(&buf, DriftDensityReport{})
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected no-resources message, got: %s", buf.String())
	}
}

func TestFprintDensity_ContainsTypes(t *testing.T) {
	var buf bytes.Buffer
	r := BuildDriftDensity(baseDensityReport())
	FprintDensity(&buf, r)
	out := buf.String()
	for _, typ := range []string{"aws_s3_bucket", "aws_instance", "aws_lambda_function"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected output to contain %q", typ)
		}
	}
	if !strings.Contains(out, "Overall drift density") {
		t.Error("expected overall density line")
	}
}
