package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseCoverageReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "sg-1", Type: "aws_security_group"},
			{ID: "sg-2", Type: "aws_security_group"},
			{ID: "s3-1", Type: "aws_s3_bucket"},
		},
		Missing: []Resource{
			{ID: "sg-3", Type: "aws_security_group"},
		},
		Untracked: []Resource{
			{ID: "rds-1", Type: "aws_db_instance"},
		},
	}
}

func TestEvaluateCoverage_Empty(t *testing.T) {
	cr := EvaluateCoverage(Report{})
	if cr.Total != 0 {
		t.Errorf("expected total 0, got %d", cr.Total)
	}
	if cr.Overall != 0.0 {
		t.Errorf("expected overall 0.0, got %f", cr.Overall)
	}
}

func TestEvaluateCoverage_CorrectCounts(t *testing.T) {
	cr := EvaluateCoverage(baseCoverageReport())
	if cr.Total != 5 {
		t.Errorf("expected total 5, got %d", cr.Total)
	}
	if cr.Managed != 3 {
		t.Errorf("expected managed 3, got %d", cr.Managed)
	}
	if cr.Overall < 59.9 || cr.Overall > 60.1 {
		t.Errorf("expected overall ~60%%, got %.2f", cr.Overall)
	}
}

func TestEvaluateCoverage_PerType(t *testing.T) {
	cr := EvaluateCoverage(baseCoverageReport())
	byType := map[string]CoverageResult{}
	for _, r := range cr.Results {
		byType[r.Type] = r
	}

	sg := byType["aws_security_group"]
	if sg.Total != 3 || sg.Managed != 2 {
		t.Errorf("sg: expected 2/3, got %d/%d", sg.Managed, sg.Total)
	}

	s3 := byType["aws_s3_bucket"]
	if s3.Total != 1 || s3.Managed != 1 {
		t.Errorf("s3: expected 1/1, got %d/%d", s3.Managed, s3.Total)
	}

	rds := byType["aws_db_instance"]
	if rds.Total != 1 || rds.Managed != 0 {
		t.Errorf("rds: expected 0/1, got %d/%d", rds.Managed, rds.Total)
	}
}

func TestCoverageHasGaps_True(t *testing.T) {
	cr := EvaluateCoverage(baseCoverageReport())
	if !CoverageHasGaps(cr) {
		t.Error("expected gaps to be detected")
	}
}

func TestCoverageHasGaps_False(t *testing.T) {
	r := Report{
		Managed: []Resource{
			{ID: "a", Type: "aws_s3_bucket"},
		},
	}
	cr := EvaluateCoverage(r)
	if CoverageHasGaps(cr) {
		t.Error("expected no gaps when all resources are managed")
	}
}

func TestFprintCoverage_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintCoverage(&buf, CoverageReport{})
	if !strings.Contains(buf.String(), "no resources") {
		t.Errorf("expected 'no resources' message, got: %s", buf.String())
	}
}

func TestFprintCoverage_OutputContainsTypes(t *testing.T) {
	cr := EvaluateCoverage(baseCoverageReport())
	var buf bytes.Buffer
	FprintCoverage(&buf, cr)
	out := buf.String()
	for _, typ := range []string{"aws_security_group", "aws_s3_bucket", "aws_db_instance"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected output to contain %q", typ)
		}
	}
	if !strings.Contains(out, "Overall") {
		t.Error("expected output to contain 'Overall'")
	}
}
