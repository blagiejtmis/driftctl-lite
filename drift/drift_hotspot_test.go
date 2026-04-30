package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseHotspotReports() []Report {
	return []Report{
		{
			Managed: []Resource{
				{ID: "vpc-1", Type: "aws_vpc"},
				{ID: "sg-1", Type: "aws_security_group"},
			},
			Missing: []Resource{
				{ID: "sg-2", Type: "aws_security_group"},
				{ID: "sg-3", Type: "aws_security_group"},
			},
			Untracked: []Resource{
				{ID: "vpc-2", Type: "aws_vpc"},
			},
		},
		{
			Managed: []Resource{
				{ID: "s3-1", Type: "aws_s3_bucket"},
			},
			Missing: []Resource{},
			Untracked: []Resource{},
		},
	}
}

func TestBuildHotspot_Empty(t *testing.T) {
	r := BuildHotspot([]Report{})
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestBuildHotspot_CorrectCounts(t *testing.T) {
	r := BuildHotspot(baseHotspotReports())

	byType := map[string]HotspotEntry{}
	for _, e := range r.Entries {
		byType[e.Type] = e
	}

	sg := byType["aws_security_group"]
	if sg.DriftCount != 2 {
		t.Errorf("expected sg DriftCount=2, got %d", sg.DriftCount)
	}
	if sg.Total != 3 {
		t.Errorf("expected sg Total=3, got %d", sg.Total)
	}

	vpc := byType["aws_vpc"]
	if vpc.DriftCount != 1 {
		t.Errorf("expected vpc DriftCount=1, got %d", vpc.DriftCount)
	}
	if vpc.Total != 2 {
		t.Errorf("expected vpc Total=2, got %d", vpc.Total)
	}

	s3 := byType["aws_s3_bucket"]
	if s3.DriftCount != 0 {
		t.Errorf("expected s3 DriftCount=0, got %d", s3.DriftCount)
	}
}

func TestBuildHotspot_SortedByScoreDesc(t *testing.T) {
	r := BuildHotspot(baseHotspotReports())
	for i := 1; i < len(r.Entries); i++ {
		if r.Entries[i].Score > r.Entries[i-1].Score {
			t.Errorf("entries not sorted by score desc at index %d", i)
		}
	}
}

func TestHotspotHasEntries_True(t *testing.T) {
	r := BuildHotspot(baseHotspotReports())
	if !HotspotHasEntries(r) {
		t.Error("expected HotspotHasEntries to be true")
	}
}

func TestHotspotHasEntries_False(t *testing.T) {
	r := BuildHotspot([]Report{
		{Managed: []Resource{{ID: "a", Type: "aws_vpc"}}},
	})
	if HotspotHasEntries(r) {
		t.Error("expected HotspotHasEntries to be false")
	}
}

func TestFprintHotspot_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintHotspot(&buf, HotspotReport{})
	if !strings.Contains(buf.String(), "No hotspot") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintHotspot_OutputContainsTypes(t *testing.T) {
	r := BuildHotspot(baseHotspotReports())
	var buf bytes.Buffer
	FprintHotspot(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "aws_security_group") {
		t.Errorf("expected aws_security_group in output, got: %s", out)
	}
	if !strings.Contains(out, "aws_vpc") {
		t.Errorf("expected aws_vpc in output, got: %s", out)
	}
	if strings.Contains(out, "aws_s3_bucket") {
		t.Errorf("s3 has no drift and should be omitted, got: %s", out)
	}
}
