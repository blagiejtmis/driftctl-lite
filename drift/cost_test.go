package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseCostReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "i-123", Type: "aws_instance", Attributes: map[string]string{"state": "running"}},
			{ID: "bucket-1", Type: "aws_s3_bucket", Attributes: map[string]string{}},
		},
		Untracked: []Resource{
			{ID: "rds-1", Type: "aws_rds_instance", Attributes: map[string]string{}},
		},
		Missing: []Resource{},
	}
}

func TestEstimateCosts_Empty(t *testing.T) {
	cr := EstimateCosts(Report{})
	if len(cr.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(cr.Entries))
	}
	if cr.TotalManaged != 0 || cr.TotalUntracked != 0 {
		t.Fatal("expected zero totals")
	}
}

func TestEstimateCosts_KnownTypes(t *testing.T) {
	cr := EstimateCosts(baseCostReport())
	if cr.TotalManaged != 52.0 { // 50 + 2
		t.Fatalf("expected TotalManaged=52.0, got %.2f", cr.TotalManaged)
	}
	if cr.TotalUntracked != 100.0 {
		t.Fatalf("expected TotalUntracked=100.0, got %.2f", cr.TotalUntracked)
	}
}

func TestEstimateCosts_UnknownType(t *testing.T) {
	r := Report{
		Managed: []Resource{
			{ID: "x-1", Type: "aws_unknown_thing", Attributes: map[string]string{}},
		},
	}
	cr := EstimateCosts(r)
	if cr.TotalManaged != 1.0 {
		t.Fatalf("expected default cost 1.0, got %.2f", cr.TotalManaged)
	}
}

func TestEstimateCosts_SortedDescending(t *testing.T) {
	cr := EstimateCosts(baseCostReport())
	for i := 1; i < len(cr.Entries); i++ {
		if cr.Entries[i].EstimatedUSD > cr.Entries[i-1].EstimatedUSD {
			t.Fatal("entries not sorted descending by cost")
		}
	}
}

func TestFprintCost_NoResources(t *testing.T) {
	var buf bytes.Buffer
	FprintCost(&buf, CostReport{})
	if !strings.Contains(buf.String(), "No resources") {
		t.Fatal("expected 'No resources' message")
	}
}

func TestFprintCost_WithResources(t *testing.T) {
	cr := EstimateCosts(baseCostReport())
	var buf bytes.Buffer
	FprintCost(&buf, cr)
	out := buf.String()
	if !strings.Contains(out, "Managed total") {
		t.Fatal("expected 'Managed total' in output")
	}
	if !strings.Contains(out, "Untracked total") {
		t.Fatal("expected 'Untracked total' in output")
	}
	if !strings.Contains(out, "untracked") {
		t.Fatal("expected 'untracked' label for rds-1")
	}
}
