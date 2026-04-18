package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baselineReport() Report {
	return Report{
		Managed: []Resource{
			{Type: "aws_s3_bucket", ID: "my-bucket", Attributes: map[string]interface{}{"region": "us-east-1"}},
			{Type: "aws_instance", ID: "i-123", Attributes: map[string]interface{}{"ami": "ami-abc"}},
		},
	}
}

func TestSaveBaseline_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := SaveBaseline(path, baselineReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	report := baselineReport()
	if err := SaveBaseline(path, report); err != nil {
		t.Fatalf("save: %v", err)
	}

	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(b.Resources) != len(report.Managed) {
		t.Errorf("expected %d resources, got %d", len(report.Managed), len(b.Resources))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if b.CreatedAt.After(time.Now().Add(5 * time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0644)

	_, err := LoadBaseline(path)
	if err == nil {
		t.Error("expected parse error")
	}
}
