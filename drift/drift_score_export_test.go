package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseScoreExport() ScoreResult {
	return ScoreResult{
		Total:     10,
		Managed:   7,
		Missing:   2,
		Untracked: 1,
		ScorePct:  70.0,
		Grade:     "C",
	}
}

func TestExportScore_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "score.json")
	result := baseScoreExport()
	if err := ExportScore(result, "json", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var got ScoreResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.Grade != "C" {
		t.Errorf("expected grade C, got %s", got.Grade)
	}
	if got.Total != 10 {
		t.Errorf("expected total 10, got %d", got.Total)
	}
}

func TestExportScore_CSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "score.csv")
	result := baseScoreExport()
	if err := ExportScore(result, "csv", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[1], "70.00") {
		t.Errorf("expected score_pct 70.00 in csv row, got: %s", lines[1])
	}
	if !strings.Contains(lines[1], "C") {
		t.Errorf("expected grade C in csv row, got: %s", lines[1])
	}
}

func TestExportScore_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "score.xml")
	err := ExportScore(baseScoreExport(), "xml", path)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExportScore_BadPath(t *testing.T) {
	err := ExportScore(baseScoreExport(), "json", "/nonexistent/dir/score.json")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}
