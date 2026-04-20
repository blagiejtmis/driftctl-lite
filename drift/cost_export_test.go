package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseCostExportReport() CostReport {
	return EstimateCosts(baseCostReport())
}

func TestExportCost_JSON(t *testing.T) {
	cr := baseCostExportReport()
	tmp := filepath.Join(t.TempDir(), "cost.json")
	if err := ExportCost(cr, "json", tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	var decoded CostReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(decoded.Entries) != len(cr.Entries) {
		t.Fatalf("entry count mismatch: want %d got %d", len(cr.Entries), len(decoded.Entries))
	}
}

func TestExportCost_CSV(t *testing.T) {
	cr := baseCostExportReport()
	tmp := filepath.Join(t.TempDir(), "cost.csv")
	if err := ExportCost(cr, "csv", tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	// header + entries
	if len(lines) != len(cr.Entries)+1 {
		t.Fatalf("expected %d lines, got %d", len(cr.Entries)+1, len(lines))
	}
	if !strings.HasPrefix(lines[0], "resource_id") {
		t.Fatal("expected csv header")
	}
}

func TestExportCost_UnsupportedFormat(t *testing.T) {
	cr := baseCostExportReport()
	tmp := filepath.Join(t.TempDir(), "cost.xml")
	err := ExportCost(cr, "xml", tmp)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestExportCost_BadPath(t *testing.T) {
	cr := baseCostExportReport()
	err := ExportCost(cr, "json", "/nonexistent/dir/cost.json")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}
