package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseAnnotations() AnnotationMap {
	return AnnotationMap{
		"aws_s3_bucket/my-bucket": {
			Type:   "aws_s3_bucket",
			ID:     "my-bucket",
			Note:   "legacy bucket, migration pending",
			Author: "alice",
		},
		"aws_iam_role/deploy-role": {
			Type: "aws_iam_role",
			ID:   "deploy-role",
			Note: "managed by ops team",
		},
	}
}

func TestAnnotateReport_NoAnnotations(t *testing.T) {
	report := Report{
		Missing: []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}},
	}
	result := AnnotateReport(report, AnnotationMap{})
	if _, ok := result.Missing[0].Attributes["_annotation"]; ok {
		t.Error("expected no annotation attribute")
	}
}

func TestAnnotateReport_MissingResource(t *testing.T) {
	report := Report{
		Missing: []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}},
	}
	result := AnnotateReport(report, baseAnnotations())
	if result.Missing[0].Attributes["_annotation"] != "legacy bucket, migration pending" {
		t.Errorf("unexpected annotation: %v", result.Missing[0].Attributes["_annotation"])
	}
}

func TestAnnotateReport_UntrackedResource(t *testing.T) {
	report := Report{
		Untracked: []Resource{{Type: "aws_iam_role", ID: "deploy-role"}},
	}
	result := AnnotateReport(report, baseAnnotations())
	if result.Untracked[0].Attributes["_annotation"] != "managed by ops team" {
		t.Errorf("unexpected annotation: %v", result.Untracked[0].Attributes["_annotation"])
	}
}

func TestFprintAnnotations_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintAnnotations(&buf, AnnotationMap{})
	if !strings.Contains(buf.String(), "No annotations") {
		t.Errorf("expected 'No annotations', got: %s", buf.String())
	}
}

func TestFprintAnnotations_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	FprintAnnotations(&buf, baseAnnotations())
	out := buf.String()
	if !strings.Contains(out, "legacy bucket") {
		t.Errorf("expected note in output, got: %s", out)
	}
	if !strings.Contains(out, "by alice") {
		t.Errorf("expected author in output, got: %s", out)
	}
	if !strings.Contains(out, "managed by ops team") {
		t.Errorf("expected second note in output, got: %s", out)
	}
}

func TestAnnotationsHasEntries(t *testing.T) {
	if AnnotationsHasEntries(AnnotationMap{}) {
		t.Error("expected false for empty map")
	}
	if !AnnotationsHasEntries(baseAnnotations()) {
		t.Error("expected true for non-empty map")
	}
}
