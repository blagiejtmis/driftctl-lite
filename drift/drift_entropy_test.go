package drift

import (
	"bytes"
	"math"
	"testing"
)

func baseEntropyReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "a", Type: "aws_s3_bucket"},
			{ID: "b", Type: "aws_s3_bucket"},
			{ID: "c", Type: "aws_instance"},
		},
		Missing: []Resource{
			{ID: "d", Type: "aws_s3_bucket"},
		},
		Untracked: []Resource{
			{ID: "e", Type: "aws_instance"},
		},
	}
}

func TestBuildDriftEntropy_Empty(t *testing.T) {
	r := BuildDriftEntropy(Report{})
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.OverallEntropy != 0 {
		t.Fatalf("expected zero overall entropy")
	}
}

func TestBuildDriftEntropy_CorrectCounts(t *testing.T) {
	r := BuildDriftEntropy(baseEntropyReport())

	found := map[string]EntropyEntry{}
	for _, e := range r.Entries {
		found[e.Type] = e
	}

	s3 := found["aws_s3_bucket"]
	if s3.Total != 3 {
		t.Errorf("s3 total: want 3, got %d", s3.Total)
	}
	if s3.Drifted != 1 {
		t.Errorf("s3 drifted: want 1, got %d", s3.Drifted)
	}

	ec2 := found["aws_instance"]
	if ec2.Total != 2 {
		t.Errorf("ec2 total: want 2, got %d", ec2.Total)
	}
	if ec2.Drifted != 1 {
		t.Errorf("ec2 drifted: want 1, got %d", ec2.Drifted)
	}
}

func TestBuildDriftEntropy_EntropyValues(t *testing.T) {
	r := BuildDriftEntropy(baseEntropyReport())
	for _, e := range r.Entries {
		if e.Entropy < 0 {
			t.Errorf("negative entropy for %s", e.Type)
		}
		if e.Entropy > 1.0+1e-9 {
			t.Errorf("entropy > 1 for %s: %.4f", e.Type, e.Entropy)
		}
	}
}

func TestBuildDriftEntropy_SortedByEntropyDesc(t *testing.T) {
	r := BuildDriftEntropy(baseEntropyReport())
	for i := 1; i < len(r.Entries); i++ {
		if r.Entries[i].Entropy > r.Entries[i-1].Entropy+1e-9 {
			t.Errorf("entries not sorted by entropy desc at index %d", i)
		}
	}
}

func TestBuildDriftEntropy_NoPureManaged(t *testing.T) {
	// A type with all managed (no drift) should have entropy 0.
	rep := Report{
		Managed: []Resource{{ID: "x", Type: "aws_vpc"}, {ID: "y", Type: "aws_vpc"}},
	}
	r := BuildDriftEntropy(rep)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if math.Abs(r.Entries[0].Entropy) > 1e-9 {
		t.Errorf("expected zero entropy, got %.4f", r.Entries[0].Entropy)
	}
}

func TestEntropyHasDrift_True(t *testing.T) {
	r := BuildDriftEntropy(baseEntropyReport())
	if !EntropyHasDrift(r) {
		t.Error("expected HasDrift=true")
	}
}

func TestEntropyHasDrift_False(t *testing.T) {
	rep := Report{Managed: []Resource{{ID: "x", Type: "aws_vpc"}}}
	r := BuildDriftEntropy(rep)
	if EntropyHasDrift(r) {
		t.Error("expected HasDrift=false")
	}
}

func TestFprintEntropy_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintEntropy(&buf, EntropyReport{})
	if buf.Len() == 0 {
		t.Error("expected non-empty output for empty report")
	}
}

func TestFprintEntropy_OutputContainsTypes(t *testing.T) {
	r := BuildDriftEntropy(baseEntropyReport())
	var buf bytes.Buffer
	FprintEntropy(&buf, r)
	out := buf.String()
	for _, e := range r.Entries {
		if !bytes.Contains([]byte(out), []byte(e.Type)) {
			t.Errorf("output missing type %s", e.Type)
		}
	}
}
