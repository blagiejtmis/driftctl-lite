package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baseAuditEntry() AuditEntry {
	return AuditEntry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Command:   "scan",
		StateFile: "terraform.tfstate",
		Total:     5,
		Missing:   1,
		Untracked: 0,
		Changed:   1,
		HasDrift:  true,
	}
}

func TestAppendAudit_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	if err := AppendAudit(path, baseAuditEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestLoadAudit_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	e := baseAuditEntry()
	_ = AppendAudit(path, e)
	_ = AppendAudit(path, e)
	entries, err := LoadAudit(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Command != "scan" {
		t.Errorf("expected command scan, got %s", entries[0].Command)
	}
}

func TestLoadAudit_MissingFile(t *testing.T) {
	entries, err := LoadAudit("/nonexistent/audit.log")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries")
	}
}

func TestAuditEntryFromReport(t *testing.T) {
	r := Report{
		Managed:   []Resource{makeResource("aws_s3_bucket", "a")},
		Missing:   []Resource{makeResource("aws_s3_bucket", "b")},
		Untracked: []Resource{},
	}
	e := AuditEntryFromReport("scan", "state.tfstate", r)
	if e.Command != "scan" {
		t.Errorf("expected scan, got %s", e.Command)
	}
	if !e.HasDrift {
		t.Errorf("expected has_drift true")
	}
	if e.Missing != 1 {
		t.Errorf("expected missing=1, got %d", e.Missing)
	}
}
