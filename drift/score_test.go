package drift

import (
	"fmt"
	"strings"
	"testing"
)

func baseScoreReport(managed, missing, untracked int) Report {
	r := Report{}
	for i := 0; i < managed; i++ {
		r.Managed = append(r.Managed, Resource{Type: "aws_s3_bucket", ID: fmt.Sprintf("m%d", i)})
	}
	for i := 0; i < missing; i++ {
		r.Missing = append(r.Missing, Resource{Type: "aws_s3_bucket", ID: fmt.Sprintf("x%d", i)})
	}
	for i := 0; i < untracked; i++ {
		r.Untracked = append(r.Untracked, Resource{Type: "aws_s3_bucket", ID: fmt.Sprintf("u%d", i)})
	}
	return r
}

func TestScoreReport_Empty(t *testing.T) {
	s := ScoreReport(Report{})
	if s.Score != 100.0 || s.Grade != "A" {
		t.Fatalf("expected A/100, got %s/%.1f", s.Grade, s.Score)
	}
}

func TestScoreReport_AllManaged(t *testing.T) {
	s := ScoreReport(baseScoreReport(10, 0, 0))
	if s.Score != 100.0 || s.Grade != "A" {
		t.Fatalf("expected A/100, got %s/%.1f", s.Grade, s.Score)
	}
}

func TestScoreReport_HalfDrift(t *testing.T) {
	s := ScoreReport(baseScoreReport(5, 5, 0))
	if s.Score != 50.0 {
		t.Fatalf("expected 50.0, got %.1f", s.Score)
	}
	if s.Grade != "C" {
		t.Fatalf("expected C, got %s", s.Grade)
	}
}

func TestScoreReport_AllDrift(t *testing.T) {
	s := ScoreReport(baseScoreReport(0, 5, 5))
	if s.Score != 0.0 || s.Grade != "F" {
		t.Fatalf("expected F/0, got %s/%.1f", s.Grade, s.Score)
	}
}

func TestFprintScore(t *testing.T) {
	s := ScoreReport(baseScoreReport(8, 1, 1))
	var sb strings.Builder
	FprintScore(&sb, s)
	out := sb.String()
	if !strings.Contains(out, "Grade") {
		t.Fatal("expected Grade in output")
	}
	if !strings.Contains(out, "80.0") {
		t.Fatalf("expected 80.0 score in output, got: %s", out)
	}
}
