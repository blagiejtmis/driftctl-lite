package drift_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"driftctl-lite/drift"
)

func makeIntegrationReport() drift.Report {
	return drift.Report{
		Managed: []drift.Resource{
			{ID: "res-1", Type: "aws_s3_bucket", Attributes: map[string]string{"region": "us-east-1"}},
			{ID: "res-2", Type: "aws_instance", Attributes: map[string]string{"ami": "ami-abc"}},
		},
		Missing: []drift.Resource{
			{ID: "res-3", Type: "aws_instance", Attributes: map[string]string{}},
		},
		Untracked: []drift.Resource{},
	}
}

func TestExportScore_Integration_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	report := makeIntegrationReport()
	result := drift.ScoreReport(report)
	if err := drift.ExportScore(result, "json", path); err != nil {
		t.Fatalf("export error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var got drift.ScoreResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json decode error: %v", err)
	}
	if got.Total != 3 {
		t.Errorf("expected total=3, got %d", got.Total)
	}
	if got.Managed != 2 {
		t.Errorf("expected managed=2, got %d", got.Managed)
	}
	if got.Missing != 1 {
		t.Errorf("expected missing=1, got %d", got.Missing)
	}
}

func TestExportScore_Integration_CSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")
	report := makeIntegrationReport()
	result := drift.ScoreReport(report)
	if err := drift.ExportScore(result, "csv", path); err != nil {
		t.Fatalf("export error: %v", err)
	}
	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "grade") {
		t.Errorf("expected header row with 'grade', got: %s", content)
	}
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines in CSV, got %d", len(lines))
	}
}
