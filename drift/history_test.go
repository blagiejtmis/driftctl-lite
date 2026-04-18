package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baseHistoryEntry() HistoryEntry {
	return HistoryEntry{
		Timestamp:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		ManagedCount:   8,
		MissingCount:   2,
		UntrackedCount: 1,
		DriftScore:     72.5,
		Label:          "ci",
	}
}

func TestAppendHistory_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.jsonl")
	err := AppendHistory(path, baseHistoryEntry())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("file not created")
	}
}

func TestLoadHistory_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.jsonl")
	e1 := baseHistoryEntry()
	e2 := baseHistoryEntry()
	e2.Label = "manual"
	e2.DriftScore = 50.0
	_ = AppendHistory(path, e1)
	_ = AppendHistory(path, e2)
	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Label != "ci" {
		t.Errorf("expected label 'ci', got %s", entries[0].Label)
	}
	if entries[1].DriftScore != 50.0 {
		t.Errorf("expected score 50.0, got %f", entries[1].DriftScore)
	}
}

func TestLoadHistory_MissingFile(t *testing.T) {
	entries, err := LoadHistory("/nonexistent/path/history.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Error("expected nil entries for missing file")
	}
}

func TestHistoryEntryFromScore(t *testing.T) {
	s := ScoreResult{Managed: 10, Missing: 2, Untracked: 1, Score: 76.9}
	e := HistoryEntryFromScore(s, "test")
	if e.ManagedCount != 10 {
		t.Errorf("expected managed 10, got %d", e.ManagedCount)
	}
	if e.Label != "test" {
		t.Errorf("expected label 'test', got %s", e.Label)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
