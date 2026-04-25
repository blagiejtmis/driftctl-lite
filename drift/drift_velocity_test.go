package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseVelocityReports() (current, previous Report) {
	current = Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc"},
			{ID: "sg-1", Type: "aws_security_group"},
		},
		Missing: []Resource{
			{ID: "vpc-2", Type: "aws_vpc"},
			{ID: "sg-2", Type: "aws_security_group"},
			{ID: "sg-3", Type: "aws_security_group"},
		},
		Untracked: []Resource{
			{ID: "s3-1", Type: "aws_s3_bucket"},
		},
	}
	previous = Report{
		Missing: []Resource{
			{ID: "vpc-2", Type: "aws_vpc"},
		},
		Untracked: []Resource{},
	}
	return
}

func TestEvaluateVelocity_NoResources(t *testing.T) {
	r := EvaluateVelocity(Report{}, Report{}, 7)
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.WindowDays != 7 {
		t.Errorf("expected window 7, got %d", r.WindowDays)
	}
}

func TestEvaluateVelocity_DefaultWindow(t *testing.T) {
	cur, prev := baseVelocityReports()
	r := EvaluateVelocity(cur, prev, 0)
	if r.WindowDays != 7 {
		t.Errorf("expected default window 7, got %d", r.WindowDays)
	}
}

func TestEvaluateVelocity_CorrectCounts(t *testing.T) {
	cur, prev := baseVelocityReports()
	r := EvaluateVelocity(cur, prev, 7)

	byType := make(map[string]VelocityEntry)
	for _, e := range r.Entries {
		byType[e.Type] = e
	}

	vpc, ok := byType["aws_vpc"]
	if !ok {
		t.Fatal("expected aws_vpc entry")
	}
	if vpc.DriftedCount != 1 {
		t.Errorf("aws_vpc drifted: want 1, got %d", vpc.DriftedCount)
	}
	if vpc.TotalCount != 2 {
		t.Errorf("aws_vpc total: want 2, got %d", vpc.TotalCount)
	}

	sg, ok := byType["aws_security_group"]
	if !ok {
		t.Fatal("expected aws_security_group entry")
	}
	if sg.DriftedCount != 2 {
		t.Errorf("aws_security_group drifted: want 2, got %d", sg.DriftedCount)
	}
}

func TestEvaluateVelocity_TrendIncreasing(t *testing.T) {
	cur, prev := baseVelocityReports()
	r := EvaluateVelocity(cur, prev, 7)

	for _, e := range r.Entries {
		if e.Type == "aws_security_group" && e.Trend != "increasing" {
			t.Errorf("aws_security_group trend: want increasing, got %s", e.Trend)
		}
		if e.Type == "aws_vpc" && e.Trend != "stable" {
			t.Errorf("aws_vpc trend: want stable, got %s", e.Trend)
		}
	}
}

func TestVelocityHasIncreasing_True(t *testing.T) {
	cur, prev := baseVelocityReports()
	r := EvaluateVelocity(cur, prev, 7)
	if !VelocityHasIncreasing(r) {
		t.Error("expected VelocityHasIncreasing to be true")
	}
}

func TestVelocityHasIncreasing_False(t *testing.T) {
	r := VelocityReport{}
	if VelocityHasIncreasing(r) {
		t.Error("expected VelocityHasIncreasing to be false")
	}
}

func TestFprintVelocity_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintVelocity(&buf, VelocityReport{WindowDays: 7})
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected 'No resources' in output, got: %s", buf.String())
	}
}

func TestFprintVelocity_OutputContainsTypes(t *testing.T) {
	cur, prev := baseVelocityReports()
	r := EvaluateVelocity(cur, prev, 14)
	var buf bytes.Buffer
	FprintVelocity(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "aws_vpc") {
		t.Errorf("expected aws_vpc in output")
	}
	if !strings.Contains(out, "14") {
		t.Errorf("expected window days in output")
	}
}
