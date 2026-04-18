package drift

import (
	"bytes"
	"testing"
)

func baseTagResources() []Resource {
	return []Resource{
		{Type: "aws_instance", ID: "i-1", Attributes: map[string]interface{}{"env": "prod", "team": "ops"}},
		{Type: "aws_instance", ID: "i-2", Attributes: map[string]interface{}{"env": "dev"}},
		{Type: "aws_s3_bucket", ID: "bucket-1", Attributes: map[string]interface{}{}},
	}
}

func TestEvaluateTags_NoViolations(t *testing.T) {
	policy := TagPolicy{Required: []TagRule{{Key: "env"}}}
	resources := []Resource{
		{Type: "aws_instance", ID: "i-1", Attributes: map[string]interface{}{"env": "prod"}},
	}
	v := EvaluateTags(resources, policy)
	if len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestEvaluateTags_MissingTag(t *testing.T) {
	policy := TagPolicy{Required: []TagRule{{Key: "team"}}}
	v := EvaluateTags(baseTagResources(), policy)
	// i-2 and bucket-1 are missing "team"
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestEvaluateTags_DisallowedValue(t *testing.T) {
	policy := TagPolicy{Required: []TagRule{{Key: "env", Values: []string{"prod"}}}}
	v := EvaluateTags(baseTagResources(), policy)
	// i-2 has env=dev (not allowed), bucket-1 missing env
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestFprintTagViolations_None(t *testing.T) {
	var buf bytes.Buffer
	FprintTagViolations(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("compliant")) {
		t.Error("expected compliant message")
	}
}

func TestFprintTagViolations_Some(t *testing.T) {
	violations := []TagViolation{
		{Resource: Resource{Type: "aws_instance", ID: "i-2"}, Rule: TagRule{Key: "team"}, Reason: "missing required tag \"team\""},
	}
	var buf bytes.Buffer
	FprintTagViolations(&buf, violations)
	if !bytes.Contains(buf.Bytes(), []byte("i-2")) {
		t.Error("expected resource id in output")
	}
}

func TestTagHasViolations(t *testing.T) {
	if TagHasViolations(nil) {
		t.Error("expected false for nil")
	}
	if !TagHasViolations([]TagViolation{{}}) {
		t.Error("expected true for non-empty")
	}
}
