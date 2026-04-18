package drift

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baseTrendEntry(score float64, grade string) TrendEntry {
	return TrendEntry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Score:     score,
		Grade:     grade,
		Managed:   10,
		Missing:   1,
		Untracked: 1,
	}
}

func TestAppendTrend_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "trend.json")
	entry := baseTrendEntry(80.0, "B")
	if err := AppendTrend(p, entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal("file not created")
	}
}

func TestLoadTrend_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "trend.json")
	e1 := baseTrendEntry(70.0, "C")
	e2 := baseTrendEntry(85.0, "B")
	_ = AppendTrend(p, e1)
	_ = AppendTrend(p, e2)
	log, err := LoadTrend(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
}

func TestLoadTrend_MissingFile(t *testing.T) {
	_, err := LoadTrend("/nonexistent/trend.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFprintTrend_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintTrend(&buf, TrendLog{})
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}

func TestFprintTrend_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	log := TrendLog{Entries: []TrendEntry{baseTrendEntry(90.0, "A")}}
	FprintTrend(&buf, log)
	if !bytes.Contains(buf.Bytes(), []byte("90.0%")) {
		t.Fatal("expected score in output")
	}
}

func TestTrendImproving(t *testing.T) {
	log := TrendLog{Entries: []TrendEntry{
		baseTrendEntry(60.0, "D"),
		baseTrendEntry(75.0, "C"),
	}}
	if !TrendImproving(log) {
		t.Fatal("expected improving trend")
	}
}
