package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseDriftIndexExport() DriftIndex {
	return BuildDriftIndex(baseDriftIndexReport())
}

func TestExportDriftIndex_JSON(t *testing.T) {
	idx := baseDriftIndexExport()
	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := ExportDriftIndex(idx, "json", tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	var got DriftIndex
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got.Entries) != len(idx.Entries) {
		t.Errorf("expected %d entries, got %d", len(idx.Entries), len(got.Entries))
	}
}

func TestExportDriftIndex_CSV(t *testing.T) {
	idx := baseDriftIndexExport()
	tmp := filepath.Join(t.TempDir(), "out.csv")
	if err := ExportDriftIndex(idx, "csv", tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	// header + entries + aggregate
	expected := 1 + len(idx.Entries) + 1
	if len(lines) != expected {
		t.Errorf("expected %d lines, got %d", expected, len(lines))
	}
	if !strings.HasPrefix(lines[0], "type,") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExportDriftIndex_UnsupportedFormat(t *testing.T) {
	idx := baseDriftIndexExport()
	err := ExportDriftIndex(idx, "xml", "/tmp/x.xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExportDriftIndex_BadPath(t *testing.T) {
	idx := baseDriftIndexExport()
	err := ExportDriftIndex(idx, "json", "/no/such/dir/out.json")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}
