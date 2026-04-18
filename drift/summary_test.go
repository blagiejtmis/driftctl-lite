package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseReport() Report {
	return Report{
		Managed:   []Resource{{Type: "aws_s3_bucket", ID: "a"}},
		Missing:   []Resource{},
		Untracked: []Resource{},
		Changed:   []Resource{},
	}
}

func TestSummarize_NoDrift(t *testing.T) {
	r := baseReport()
	s := Summarize(r)
	if s.Total != 1 || s.Managed != 1 || s.HasDrift() {
		t.Errorf("expected no drift, got %+v", s)
	}
}

func TestSummarize_WithMissing(t *testing.T) {
	r := baseReport()
	r.Missing = append(r.Missing, Resource{Type: "aws_instance", ID: "i-123"})
	s := Summarize(r)
	if s.Missing != 1 || !s.HasDrift() {
		t.Errorf("expected drift via missing, got %+v", s)
	}
	if s.Total != 2 {
		t.Errorf("expected total 2, got %d", s.Total)
	}
}

func TestSummarize_WithUntracked(t *testing.T) {
	r := baseReport()
	r.Untracked = append(r.Untracked, Resource{Type: "aws_sg", ID: "sg-1"})
	s := Summarize(r)
	if !s.HasDrift() || s.Untracked != 1 {
		t.Errorf("expected drift via untracked, got %+v", s)
	}
}

func TestFprintSummary_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	FprintSummary(&buf, Summary{Total: 3, Managed: 3})
	out := buf.String()
	if !strings.Contains(out, "No drift detected.") {
		t.Errorf("expected no drift message, got: %s", out)
	}
}

func TestFprintSummary_HasDrift(t *testing.T) {
	var buf bytes.Buffer
	FprintSummary(&buf, Summary{Total: 2, Managed: 1, Missing: 1})
	out := buf.String()
	if !strings.Contains(out, "Drift detected!") {
		t.Errorf("expected drift message, got: %s", out)
	}
}
