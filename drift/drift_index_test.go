package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseDriftIndexReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "a", Type: "aws_s3_bucket"},
			{ID: "b", Type: "aws_s3_bucket"},
			{ID: "c", Type: "aws_instance"},
		},
		Missing: []Resource{
			{ID: "d", Type: "aws_s3_bucket"},
			{ID: "e", Type: "aws_rds_instance"},
		},
		Untracked: []Resource{
			{ID: "f", Type: "aws_instance"},
		},
		Changed: []Resource{
			{ID: "a", Type: "aws_s3_bucket"},
		},
	}
}

func TestBuildDriftIndex_Empty(t *testing.T) {
	idx := BuildDriftIndex(Report{})
	if len(idx.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(idx.Entries))
	}
	if idx.Total.Total != 0 {
		t.Fatalf("expected zero total")
	}
}

func TestBuildDriftIndex_CorrectCounts(t *testing.T) {
	idx := BuildDriftIndex(baseDriftIndexReport())

	var s3 *IndexEntry
	for i := range idx.Entries {
		if idx.Entries[i].Type == "aws_s3_bucket" {
			s3 = &idx.Entries[i]
		}
	}
	if s3 == nil {
		t.Fatal("expected aws_s3_bucket entry")
	}
	if s3.Total != 3 {
		t.Errorf("expected 3 total for s3, got %d", s3.Total)
	}
	if s3.Missing != 1 {
		t.Errorf("expected 1 missing for s3, got %d", s3.Missing)
	}
	if s3.Changed != 1 {
		t.Errorf("expected 1 changed for s3, got %d", s3.Changed)
	}
}

func TestBuildDriftIndex_AggregateTotal(t *testing.T) {
	idx := BuildDriftIndex(baseDriftIndexReport())
	if idx.Total.Total != 6 {
		t.Errorf("expected total 6, got %d", idx.Total.Total)
	}
	if idx.Total.Missing != 2 {
		t.Errorf("expected missing 2, got %d", idx.Total.Missing)
	}
	if idx.Total.Untracked != 1 {
		t.Errorf("expected untracked 1, got %d", idx.Total.Untracked)
	}
}

func TestBuildDriftIndex_SortedByDriftPctDesc(t *testing.T) {
	idx := BuildDriftIndex(baseDriftIndexReport())
	for i := 1; i < len(idx.Entries); i++ {
		if idx.Entries[i].DriftPct > idx.Entries[i-1].DriftPct {
			t.Errorf("entries not sorted descending by DriftPct at index %d", i)
		}
	}
}

func TestDriftIndexHasDrift_True(t *testing.T) {
	idx := BuildDriftIndex(baseDriftIndexReport())
	if !DriftIndexHasDrift(idx) {
		t.Error("expected drift")
	}
}

func TestDriftIndexHasDrift_False(t *testing.T) {
	r := Report{
		Managed: []Resource{{ID: "a", Type: "aws_s3_bucket"}},
	}
	idx := BuildDriftIndex(r)
	if DriftIndexHasDrift(idx) {
		t.Error("expected no drift")
	}
}

func TestFprintDriftIndex_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintDriftIndex(&buf, DriftIndex{})
	if !strings.Contains(buf.String(), "no resources") {
		t.Errorf("expected 'no resources' message, got: %s", buf.String())
	}
}

func TestFprintDriftIndex_ContainsTypes(t *testing.T) {
	idx := BuildDriftIndex(baseDriftIndexReport())
	var buf bytes.Buffer
	FprintDriftIndex(&buf, idx)
	out := buf.String()
	for _, typ := range []string{"aws_s3_bucket", "aws_instance", "aws_rds_instance", "(all)"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected %q in output", typ)
		}
	}
}
