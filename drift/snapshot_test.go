package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func baseSnapshotReport() Report {
	return Report{
		Managed: []Resource{
			{Type: "aws_s3_bucket", ID: "bucket-1", Attributes: map[string]string{"region": "us-east-1"}},
		},
		Missing:   []Resource{},
		Untracked: []Resource{},
	}
}

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	report := baseSnapshotReport()

	if err := SaveSnapshot(path, "test-label", report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	report := baseSnapshotReport()

	if err := SaveSnapshot(path, "round-trip", report); err != nil {
		t.Fatalf("save: %v", err)
	}
	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if snap.Label != "round-trip" {
		t.Errorf("expected label 'round-trip', got %q", snap.Label)
	}
	if len(snap.Report.Managed) != 1 {
		t.Errorf("expected 1 managed resource, got %d", len(snap.Report.Managed))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected parse error")
	}
}
