package drift

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseAgingReport() Report {
	return Report{
		Managed:   []Resource{{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{}}},
		Missing:   []Resource{{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{}}},
		Untracked: []Resource{{ID: "s3-1", Type: "aws_s3_bucket", Attributes: map[string]string{}}},
	}
}

func TestEvaluateAging_NoResources(t *testing.T) {
	report := Report{}
	ar := EvaluateAging(report, nil)
	if AgingHasEntries(ar) {
		t.Fatal("expected no aging entries for empty report")
	}
}

func TestEvaluateAging_EntriesCreated(t *testing.T) {
	report := baseAgingReport()
	ar := EvaluateAging(report, nil)
	if len(ar.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ar.Entries))
	}
}

func TestEvaluateAging_AgeFromAudit(t *testing.T) {
	report := baseAgingReport()
	oldTime := time.Now().UTC().Add(-72 * time.Hour)
	audit := []AuditEntry{
		{
			Timestamp: oldTime,
			Missing:   []Resource{{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{}}},
		},
	}
	ar := EvaluateAging(report, audit)
	var found *AgingEntry
	for i := range ar.Entries {
		if ar.Entries[i].ResourceID == "sg-1" {
			found = &ar.Entries[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected sg-1 in aging entries")
	}
	if found.AgeDays < 2 {
		t.Errorf("expected age >= 2 days, got %d", found.AgeDays)
	}
}

func TestEvaluateAging_SortedDescending(t *testing.T) {
	report := baseAgingReport()
	now := time.Now().UTC()
	audit := []AuditEntry{
		{Timestamp: now.Add(-24 * time.Hour), Missing: []Resource{{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{}}}},
		{Timestamp: now.Add(-96 * time.Hour), Untracked: []Resource{{ID: "s3-1", Type: "aws_s3_bucket", Attributes: map[string]string{}}}},
	}
	ar := EvaluateAging(report, audit)
	if len(ar.Entries) < 2 {
		t.Fatal("expected at least 2 entries")
	}
	if ar.Entries[0].AgeDays < ar.Entries[1].AgeDays {
		t.Error("entries should be sorted by age descending")
	}
}

func TestFprintAging_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintAging(&buf, AgingReport{})
	if !strings.Contains(buf.String(), "no drifted") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestFprintAging_WithEntries(t *testing.T) {
	report := baseAgingReport()
	ar := EvaluateAging(report, nil)
	var buf bytes.Buffer
	FprintAging(&buf, ar)
	out := buf.String()
	if !strings.Contains(out, "sg-1") || !strings.Contains(out, "s3-1") {
		t.Errorf("expected resource IDs in output, got: %s", out)
	}
}
