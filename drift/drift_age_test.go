package drift

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseDriftAgeReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "res-1", Type: "aws_s3_bucket"},
		},
		Missing: []Resource{
			{ID: "res-2", Type: "aws_instance"},
		},
		Untracked: []Resource{
			{ID: "res-3", Type: "aws_security_group"},
		},
	}
}

func TestEvaluateDriftAge_NoResources(t *testing.T) {
	report := Report{}
	result := EvaluateDriftAge(report, nil)
	if DriftAgeHasEntries(result) {
		t.Errorf("expected no entries, got %d", len(result.Entries))
	}
}

func TestEvaluateDriftAge_EntriesCreated(t *testing.T) {
	report := baseDriftAgeReport()
	result := EvaluateDriftAge(report, nil)
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestEvaluateDriftAge_KindsCorrect(t *testing.T) {
	report := baseDriftAgeReport()
	result := EvaluateDriftAge(report, nil)
	kinds := map[string]bool{}
	for _, e := range result.Entries {
		kinds[e.DriftKind] = true
	}
	if !kinds["missing"] || !kinds["untracked"] {
		t.Errorf("expected both 'missing' and 'untracked' kinds, got %v", kinds)
	}
}

func TestEvaluateDriftAge_UsesAuditFirstSeen(t *testing.T) {
	report := baseDriftAgeReport()
	old := time.Now().UTC().Add(-72 * time.Hour)
	audit := []AuditEntry{
		{
			Timestamp: old,
			Missing:   []Resource{{ID: "res-2", Type: "aws_instance"}},
		},
	}
	result := EvaluateDriftAge(report, audit)
	for _, e := range result.Entries {
		if e.ResourceID == "res-2" {
			if e.AgeDays < 2.9 {
				t.Errorf("expected age >= 3 days, got %.2f", e.AgeDays)
			}
			return
		}
	}
	t.Error("entry for res-2 not found")
}

func TestEvaluateDriftAge_SortedDescending(t *testing.T) {
	report := baseDriftAgeReport()
	old := time.Now().UTC().Add(-48 * time.Hour)
	audit := []AuditEntry{
		{Timestamp: old, Missing: []Resource{{ID: "res-2", Type: "aws_instance"}}},
	}
	result := EvaluateDriftAge(report, audit)
	for i := 1; i < len(result.Entries); i++ {
		if result.Entries[i].AgeDays > result.Entries[i-1].AgeDays {
			t.Error("entries not sorted descending by age")
		}
	}
}

func TestFprintDriftAge_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintDriftAge(&buf, DriftAgeReport{})
	if !strings.Contains(buf.String(), "No drift-age") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintDriftAge_OutputContainsIDs(t *testing.T) {
	report := baseDriftAgeReport()
	result := EvaluateDriftAge(report, nil)
	var buf bytes.Buffer
	FprintDriftAge(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "res-2") || !strings.Contains(out, "res-3") {
		t.Errorf("expected resource IDs in output, got: %s", out)
	}
}
