package drift

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempState(t *testing.T, state tfState) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.tfstate")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(state); err != nil {
		t.Fatalf("encode state: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestScan_EmptyState(t *testing.T) {
	path := writeTempState(t, tfState{Version: 4, Resources: nil})
	result, err := Scan(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Managed)+len(result.Drifted)+len(result.Missing) != 0 {
		t.Errorf("expected empty result, got %+v", result)
	}
}

func TestScan_MissingFile(t *testing.T) {
	_, err := Scan("/nonexistent/path.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseResources(t *testing.T) {
	state := tfState{
		Version: 4,
		Resources: []tfResource{
			{Type: "aws_s3_bucket", Name: "my_bucket", Mode: "managed",
				Instances: []tfInstance{{Attributes: map[string]interface{}{"bucket": "my-bucket"}}}},
		},
	}
	resources := parseResources(state)
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	if resources[0].Type != "aws_s3_bucket" {
		t.Errorf("unexpected type %q", resources[0].Type)
	}
}
