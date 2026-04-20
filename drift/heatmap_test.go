package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseHeatmapReport() Report {
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
			{ID: "f", Type: "aws_lambda_function"},
		},
	}
}

func TestBuildHeatmap_Empty(t *testing.T) {
	h := BuildHeatmap(Report{})
	if len(h.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(h.Entries))
	}
}

func TestBuildHeatmap_CorrectIntensity(t *testing.T) {
	h := BuildHeatmap(baseHeatmapReport())

	byType := map[string]HeatmapEntry{}
	for _, e := range h.Entries {
		byType[e.Type] = e
	}

	s3 := byType["aws_s3_bucket"]
	if s3.Total != 3 {
		t.Errorf("s3 total: want 3, got %d", s3.Total)
	}
	if s3.Drifted != 1 {
		t.Errorf("s3 drifted: want 1, got %d", s3.Drifted)
	}
	want := 1.0 / 3.0
	if s3.Intensity < want-0.001 || s3.Intensity > want+0.001 {
		t.Errorf("s3 intensity: want ~%.3f, got %.3f", want, s3.Intensity)
	}

	lambda := byType["aws_lambda_function"]
	if lambda.Intensity != 1.0 {
		t.Errorf("lambda intensity: want 1.0, got %.3f", lambda.Intensity)
	}
}

func TestBuildHeatmap_SortedByIntensityDesc(t *testing.T) {
	h := BuildHeatmap(baseHeatmapReport())
	for i := 1; i < len(h.Entries); i++ {
		if h.Entries[i].Intensity > h.Entries[i-1].Intensity {
			t.Errorf("entries not sorted desc at index %d", i)
		}
	}
}

func TestFprintHeatmap_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintHeatmap(&buf, HeatmapResult{})
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintHeatmap_ContainsTypes(t *testing.T) {
	var buf bytes.Buffer
	FprintHeatmap(&buf, BuildHeatmap(baseHeatmapReport()))
	out := buf.String()
	for _, typ := range []string{"aws_s3_bucket", "aws_instance", "aws_lambda_function"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected %q in output", typ)
		}
	}
	if !strings.Contains(out, "100%") {
		t.Errorf("expected 100%% intensity for lambda in output")
	}
}
