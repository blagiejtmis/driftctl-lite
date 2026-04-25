package drift

import (
	"bytes"
	"strings"
	"testing"
)

func basePruneReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc"},
		},
		Missing: []Resource{
			{ID: "sg-old", Type: "aws_security_group"},
		},
		Untracked: []Resource{
			{ID: "s3-orphan", Type: "aws_s3_bucket"},
			{ID: "ec2-orphan", Type: "aws_instance"},
		},
	}
}

func TestEvaluatePrune_NoResources(t *testing.T) {
	pr, err := EvaluatePrune(Report{}, "", PruneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if PruneHasEntries(pr) {
		t.Error("expected no prune entries for empty report")
	}
}

func TestEvaluatePrune_UntrackedIncluded(t *testing.T) {
	report := basePruneReport()
	pr, err := EvaluatePrune(report, "", PruneOptions{OnlyUntracked: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pr.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(pr.Entries))
	}
}

func TestEvaluatePrune_FilterByType(t *testing.T) {
	report := basePruneReport()
	opts := PruneOptions{
		OnlyUntracked: true,
		Types:         []string{"aws_s3_bucket"},
	}
	pr, err := EvaluatePrune(report, "", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pr.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(pr.Entries))
	}
	if pr.Entries[0].Resource.ID != "s3-orphan" {
		t.Errorf("expected s3-orphan, got %s", pr.Entries[0].Resource.ID)
	}
}

func TestEvaluatePrune_MaxAgeDays_FiltersAll(t *testing.T) {
	report := basePruneReport()
	// age map will be empty (no audit file), so all age=0; MaxAgeDays=1 should exclude all
	opts := PruneOptions{OnlyUntracked: true, MaxAgeDays: 1}
	pr, err := EvaluatePrune(report, "", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if PruneHasEntries(pr) {
		t.Errorf("expected no entries when age threshold not met, got %d", len(pr.Entries))
	}
}

func TestFprintPrune_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	FprintPrune(&buf, PruneReport{})
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected no-resources message, got: %s", buf.String())
	}
}

func TestFprintPrune_WithEntries(t *testing.T) {
	pr := PruneReport{
		Evaluated: 3,
		Entries: []PruneEntry{
			{Resource: Resource{ID: "s3-orphan", Type: "aws_s3_bucket"}, Reason: "untracked", AgeDays: 0},
		},
	}
	var buf bytes.Buffer
	FprintPrune(&buf, pr)
	out := buf.String()
	if !strings.Contains(out, "s3-orphan") {
		t.Errorf("expected s3-orphan in output, got: %s", out)
	}
	if !strings.Contains(out, "1 of 3") {
		t.Errorf("expected count in output, got: %s", out)
	}
}
