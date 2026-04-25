package drift

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func baseMaturityExport() MaturityReport {
	return EvaluateMaturity(baseMaturityReport())
}

func TestExportMaturity_JSON(t *testing.T) {
	r := baseMaturityExport()
	tmp, err := os.CreateTemp("", "maturity-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := ExportMaturity(r, "json", tmp.Name()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp.Name())
	var out MaturityReport
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Results) != len(r.Results) {
		t.Errorf("expected %d results, got %d", len(r.Results), len(out.Results))
	}
}

func TestExportMaturity_CSV(t *testing.T) {
	r := baseMaturityExport()
	tmp, err := os.CreateTemp("", "maturity-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := ExportMaturity(r, "csv", tmp.Name()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp.Name())
	if !strings.Contains(string(data), "type,total,managed") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(string(data), "aws_s3_bucket") {
		t.Error("expected aws_s3_bucket in CSV")
	}
}

func TestExportMaturity_UnsupportedFormat(t *testing.T) {
	err := ExportMaturity(MaturityReport{}, "xml", "/tmp/out.xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportMaturity_BadPath(t *testing.T) {
	err := ExportMaturity(baseMaturityExport(), "json", "/nonexistent/dir/out.json")
	if err == nil {
		t.Error("expected error for bad path")
	}
}
