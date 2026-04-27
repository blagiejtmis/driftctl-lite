package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseSurfaceReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "r1", Type: "aws_s3_bucket"},
			{ID: "r2", Type: "aws_s3_bucket"},
			{ID: "r3", Type: "aws_instance"},
		},
		Missing: []Resource{
			{ID: "r4", Type: "aws_s3_bucket"},
		},
		Untracked: []Resource{
			{ID: "r5", Type: "aws_instance"},
			{ID: "r6", Type: "aws_instance"},
		},
	}
}

func TestBuildDriftSurface_Empty(t *testing.T) {
	r := BuildDriftSurface(Report{})
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.Overall != 0 {
		t.Errorf("expected 0 overall, got %f", r.Overall)
	}
}

func TestBuildDriftSurface_CorrectCounts(t *testing.T) {
	r := BuildDriftSurface(baseSurfaceReport())
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	for _, e := range r.Entries {
		switch e.Type {
		case "aws_s3_bucket":
			if e.Total != 3 || e.Drifted != 1 {
				t.Errorf("s3 bucket: want total=3 drifted=1, got %d/%d", e.Total, e.Drifted)
			}
		case "aws_instance":
			if e.Total != 3 || e.Drifted != 2 {
				t.Errorf("instance: want total=3 drifted=2, got %d/%d", e.Total, e.Drifted)
			}
		default:
			t.Errorf("unexpected type %s", e.Type)
		}
	}
}

func TestBuildDriftSurface_SortedByPctDesc(t *testing.T) {
	r := BuildDriftSurface(baseSurfaceReport())
	for i := 1; i < len(r.Entries); i++ {
		if r.Entries[i].SurfacePct > r.Entries[i-1].SurfacePct {
			t.Errorf("entries not sorted desc at index %d", i)
		}
	}
}

func TestBuildDriftSurface_OverallPct(t *testing.T) {
	r := BuildDriftSurface(baseSurfaceReport())
	// total=6, drifted=3 => 50%
	if r.Overall != 50.0 {
		t.Errorf("expected overall 50.0, got %f", r.Overall)
	}
}

func TestSurfaceHasCritical_True(t *testing.T) {
	r := SurfaceReport{
		Entries: []SurfaceEntry{{Type: "x", SurfacePct: 80, RiskLevel: "critical"}},
	}
	if !SurfaceHasCritical(r) {
		t.Error("expected critical")
	}
}

func TestSurfaceHasCritical_False(t *testing.T) {
	r := SurfaceReport{
		Entries: []SurfaceEntry{{Type: "x", SurfacePct: 10, RiskLevel: "low"}},
	}
	if SurfaceHasCritical(r) {
		t.Error("expected no critical")
	}
}

func TestFprintSurface_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintSurface(&buf, SurfaceReport{})
	if !strings.Contains(buf.String(), "no resources") {
		t.Errorf("expected 'no resources', got: %s", buf.String())
	}
}

func TestFprintSurface_OutputContainsTypes(t *testing.T) {
	r := BuildDriftSurface(baseSurfaceReport())
	var buf bytes.Buffer
	FprintSurface(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Error("expected aws_s3_bucket in output")
	}
	if !strings.Contains(out, "aws_instance") {
		t.Error("expected aws_instance in output")
	}
	if !strings.Contains(out, "overall") {
		t.Error("expected 'overall' in output")
	}
}
