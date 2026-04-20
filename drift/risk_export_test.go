package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseRiskExportResults() []RiskResult {
	return []RiskResult{
		{ResourceID: "r1", ResourceType: "aws_s3_bucket", Level: "high", Reason: "public bucket", Score: 70},
		{ResourceID: "r2", ResourceType: "aws_iam_role", Level: "critical", Reason: "admin policy", Score: 95},
	}
}

func TestExportRisk_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "risk.json")
	if err := ExportRisk(baseRiskExportResults(), "json", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var out []RiskResult
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestExportRisk_CSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "risk.csv")
	if err := ExportRisk(baseRiskExportResults(), "csv", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "resource_id") {
		t.Errorf("expected csv header, got %q", lines[0])
	}
}

func TestExportRisk_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "risk.xml")
	err := ExportRisk(baseRiskExportResults(), "xml", path)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportRisk_BadPath(t *testing.T) {
	err := ExportRisk(baseRiskExportResults(), "json", "/nonexistent/dir/risk.json")
	if err == nil {
		t.Error("expected error for bad path")
	}
}
