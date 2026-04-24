package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseForecastExport() ForecastResult {
	return ForecastResult{
		Trend:      "improving",
		Confidence: "medium",
		Entries: []ForecastEntry{
			{Period: "+1", Score: 82.5, Grade: "B"},
			{Period: "+2", Score: 85.0, Grade: "B"},
		},
	}
}

func TestExportForecast_JSON(t *testing.T) {
	r := baseForecastExport()
	path := filepath.Join(t.TempDir(), "forecast.json")
	if err := ExportForecast(r, "json", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var got ForecastResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got.Entries))
	}
}

func TestExportForecast_CSV(t *testing.T) {
	r := baseForecastExport()
	path := filepath.Join(t.TempDir(), "forecast.csv")
	if err := ExportForecast(r, "csv", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header+2), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "period") {
		t.Errorf("expected CSV header, got %s", lines[0])
	}
}

func TestExportForecast_UnsupportedFormat(t *testing.T) {
	err := ExportForecast(baseForecastExport(), "xml", "/tmp/x.xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportForecast_BadPath(t *testing.T) {
	err := ExportForecast(baseForecastExport(), "json", "/nonexistent/dir/f.json")
	if err == nil {
		t.Error("expected error for bad path")
	}
}
