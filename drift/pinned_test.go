package drift

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func basePinnedResources() []Resource {
	return []Resource{
		{Type: "aws_s3_bucket", ID: "bucket-1", Attributes: map[string]string{"region": "us-east-1"}},
		{Type: "aws_instance", ID: "i-001", Attributes: map[string]string{"ami": "ami-123"}},
	}
}

func TestPinResources_Empty(t *testing.T) {
	ps := PinResources(nil)
	if len(ps.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(ps.Entries))
	}
}

func TestPinResources_PopulatesEntries(t *testing.T) {
	ps := PinResources(basePinnedResources())
	if len(ps.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ps.Entries))
	}
	// sorted by type then id: aws_instance before aws_s3_bucket
	if ps.Entries[0].Type != "aws_instance" {
		t.Errorf("expected first entry type aws_instance, got %s", ps.Entries[0].Type)
	}
	if ps.Entries[0].Attributes["ami"] != "ami-123" {
		t.Errorf("attribute not copied correctly")
	}
}

func TestSavePinned_RoundTrip(t *testing.T) {
	ps := PinResources(basePinnedResources())
	tmp := filepath.Join(t.TempDir(), "pinned.json")
	if err := SavePinned(tmp, ps); err != nil {
		t.Fatalf("SavePinned: %v", err)
	}
	loaded, err := LoadPinned(tmp)
	if err != nil {
		t.Fatalf("LoadPinned: %v", err)
	}
	if len(loaded.Entries) != len(ps.Entries) {
		t.Errorf("entry count mismatch: want %d got %d", len(ps.Entries), len(loaded.Entries))
	}
}

func TestLoadPinned_MissingFile(t *testing.T) {
	_, err := LoadPinned(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadPinned_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not-json"), 0o644)
	_, err := LoadPinned(tmp)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSavePinned_ValidJSON(t *testing.T) {
	ps := PinResources(basePinnedResources())
	tmp := filepath.Join(t.TempDir(), "pinned.json")
	SavePinned(tmp, ps)
	data, _ := os.ReadFile(tmp)
	var check PinnedSet
	if err := json.Unmarshal(data, &check); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestFprintPinned_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintPinned(&buf, PinnedSet{})
	if buf.String() != "No pinned resources.\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFprintPinned_WithEntries(t *testing.T) {
	ps := PinResources(basePinnedResources())
	var buf bytes.Buffer
	FprintPinned(&buf, ps)
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !bytes.Contains([]byte(out), []byte("aws_instance")) {
		t.Errorf("expected aws_instance in output")
	}
}
