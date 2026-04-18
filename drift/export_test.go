package drift

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func baseExportReport() Report {
	return Report{
		Managed:   []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}},
		Missing:   []Resource{{Type: "aws_instance", ID: "i-123"}},
		Untracked: []Resource{{Type: "aws_vpc", ID: "vpc-abc"}},
	}
}

func TestExport_JSON(t *testing.T) {
	f, err := os.CreateTemp("", "export-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	report := baseExportReport()
	if err := Export(report, f.Name(), ExportJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(f.Name())
	var got Report
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got.Missing) != 1 || got.Missing[0].ID != "i-123" {
		t.Errorf("unexpected missing: %+v", got.Missing)
	}
}

func TestExport_CSV(t *testing.T) {
	f, err := os.CreateTemp("", "export-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	report := baseExportReport()
	if err := Export(report, f.Name(), ExportCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(f.Name())
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "status,type,id" {
		t.Errorf("bad header: %q", lines[0])
	}
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	err := Export(Report{}, "/tmp/x", ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExport_BadPath(t *testing.T) {
	err := Export(Report{}, "/nonexistent/dir/out.json", ExportJSON)
	if err == nil {
		t.Error("expected error for bad path")
	}
}
