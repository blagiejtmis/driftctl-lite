package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseDriftRatioReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc"},
			{ID: "sg-1", Type: "aws_security_group"},
			{ID: "sg-2", Type: "aws_security_group"},
		},
		Missing: []Resource{
			{ID: "vpc-2", Type: "aws_vpc"},
		},
		Untracked: []Resource{
			{ID: "sg-3", Type: "aws_security_group"},
		},
	}
}

func TestBuildDriftRatio_Empty(t *testing.T) {
	rr := BuildDriftRatio(Report{})
	if len(rr.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(rr.Entries))
	}
	if rr.OverallRatio != 0 {
		t.Fatalf("expected overall ratio 0, got %f", rr.OverallRatio)
	}
}

func TestBuildDriftRatio_CorrectCounts(t *testing.T) {
	rr := BuildDriftRatio(baseDriftRatioReport())

	byType := map[string]RatioEntry{}
	for _, e := range rr.Entries {
		byType[e.Type] = e
	}

	vpc := byType["aws_vpc"]
	if vpc.Total != 2 || vpc.Drifted != 1 {
		t.Errorf("aws_vpc: want total=2 drifted=1, got total=%d drifted=%d", vpc.Total, vpc.Drifted)
	}
	if vpc.Ratio != 0.5 {
		t.Errorf("aws_vpc ratio: want 0.5, got %f", vpc.Ratio)
	}

	sg := byType["aws_security_group"]
	if sg.Total != 3 || sg.Drifted != 1 {
		t.Errorf("aws_security_group: want total=3 drifted=1, got total=%d drifted=%d", sg.Total, sg.Drifted)
	}
}

func TestBuildDriftRatio_OverallRatio(t *testing.T) {
	rr := BuildDriftRatio(baseDriftRatioReport())
	// 5 total, 2 drifted → 0.4
	if rr.OverallRatio < 0.39 || rr.OverallRatio > 0.41 {
		t.Errorf("overall ratio: want ~0.4, got %f", rr.OverallRatio)
	}
}

func TestBuildDriftRatio_SortedByRatioDesc(t *testing.T) {
	rr := BuildDriftRatio(baseDriftRatioReport())
	for i := 1; i < len(rr.Entries); i++ {
		if rr.Entries[i].Ratio > rr.Entries[i-1].Ratio {
			t.Errorf("entries not sorted descending at index %d", i)
		}
	}
}

func TestDriftRatioHasDrift_True(t *testing.T) {
	rr := BuildDriftRatio(baseDriftRatioReport())
	if !DriftRatioHasDrift(rr) {
		t.Error("expected HasDrift=true")
	}
}

func TestDriftRatioHasDrift_False(t *testing.T) {
	rr := BuildDriftRatio(Report{
		Managed: []Resource{{ID: "vpc-1", Type: "aws_vpc"}},
	})
	if DriftRatioHasDrift(rr) {
		t.Error("expected HasDrift=false")
	}
}

func TestFprintDriftRatio_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintDriftRatio(&buf, RatioReport{})
	if !strings.Contains(buf.String(), "no resources") {
		t.Errorf("expected 'no resources' message, got: %s", buf.String())
	}
}

func TestFprintDriftRatio_OutputContainsTypes(t *testing.T) {
	rr := BuildDriftRatio(baseDriftRatioReport())
	var buf bytes.Buffer
	FprintDriftRatio(&buf, rr)
	out := buf.String()
	for _, want := range []string{"aws_vpc", "aws_security_group", "Overall"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}
