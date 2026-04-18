package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempPolicy(t *testing.T, pf PolicyFile) string {
	t.Helper()
	data, _ := json.Marshal(pf)
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.json")
	_ = os.WriteFile(p, data, 0644)
	return p
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	_, err := LoadPolicy("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadPolicy_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(p, []byte("not-json"), 0644)
	_, err := LoadPolicy(p)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestLoadPolicy_Valid(t *testing.T) {
	pf := PolicyFile{Rules: []PolicyRule{{ID: "r1", Type: "aws_s3_bucket", Severity: "error", Message: "missing bucket"}}}
	p := writeTempPolicy(t, pf)
	loaded, err := LoadPolicy(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(loaded.Rules))
	}
}

func TestEvaluatePolicy_NoViolations(t *testing.T) {
	pf := &PolicyFile{Rules: []PolicyRule{{ID: "r1", Type: "aws_s3_bucket", Severity: "error"}}}
	report := Report{}
	results := EvaluatePolicy(pf, report)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEvaluatePolicy_WithViolations(t *testing.T) {
	pf := &PolicyFile{Rules: []PolicyRule{{ID: "r1", Type: "aws_s3_bucket", Severity: "error", Message: "missing"}}}
	report := Report{
		Missing: []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}},
	}
	results := EvaluatePolicy(pf, report)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Violations) != 1 {
		t.Fatalf("expected 1 violation")
	}
}

func TestEvaluatePolicy_TypeMismatch(t *testing.T) {
	pf := &PolicyFile{Rules: []PolicyRule{{ID: "r1", Type: "aws_lambda_function", Severity: "warn"}}}
	report := Report{
		Missing: []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}},
	}
	results := EvaluatePolicy(pf, report)
	if len(results) != 0 {
		t.Fatalf("expected 0 results for type mismatch")
	}
}
