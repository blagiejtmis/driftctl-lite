package drift

import (
	"bytes"
	"testing"
)

var baseLabelResources = []Resource{
	{ID: "res-1", Type: "aws_instance", Attributes: map[string]interface{}{"env": "prod", "team": "platform"}},
	{ID: "res-2", Type: "aws_s3_bucket", Attributes: map[string]interface{}{"env": "staging"}},
	{ID: "res-3", Type: "aws_rds_instance", Attributes: map[string]interface{}{}},
}

var baseLabelRules = []LabelRule{
	{Key: "env", Required: true, Allowed: []string{"prod", "staging", "dev"}},
	{Key: "team", Required: true},
}

func TestEvaluateLabels_NoViolations(t *testing.T) {
	res := []Resource{
		{ID: "r1", Type: "aws_instance", Attributes: map[string]interface{}{"env": "prod", "team": "sre"}},
	}
	got := EvaluateLabels(res, baseLabelRules)
	if len(got) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(got))
	}
}

func TestEvaluateLabels_MissingLabel(t *testing.T) {
	got := EvaluateLabels(baseLabelResources, baseLabelRules)
	// res-2 missing team, res-3 missing env and team
	if len(got) < 3 {
		t.Fatalf("expected at least 3 violations, got %d", len(got))
	}
}

func TestEvaluateLabels_DisallowedValue(t *testing.T) {
	res := []Resource{
		{ID: "r1", Type: "aws_instance", Attributes: map[string]interface{}{"env": "unknown", "team": "ops"}},
	}
	got := EvaluateLabels(res, baseLabelRules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Rule.Key != "env" {
		t.Errorf("expected violation on 'env', got %q", got[0].Rule.Key)
	}
}

func TestEvaluateLabels_SortedOutput(t *testing.T) {
	got := EvaluateLabels(baseLabelResources, baseLabelRules)
	for i := 1; i < len(got); i++ {
		ki := resourceKey(got[i-1].Resource)
		kj := resourceKey(got[i].Resource)
		if ki > kj {
			t.Errorf("violations not sorted: %s > %s", ki, kj)
		}
	}
}

func TestLabelHasViolations_True(t *testing.T) {
	v := []LabelViolation{{Message: "x"}}
	if !LabelHasViolations(v) {
		t.Error("expected true")
	}
}

func TestLabelHasViolations_False(t *testing.T) {
	if LabelHasViolations(nil) {
		t.Error("expected false")
	}
}

func TestFprintLabels_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	FprintLabels(&buf, nil)
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("no violations")) {
		t.Error("expected 'no violations' in output")
	}
}

func TestFprintLabels_WithViolations(t *testing.T) {
	v := []LabelViolation{
		{Resource: Resource{ID: "r1", Type: "aws_instance"}, Rule: LabelRule{Key: "env"}, Message: "missing required label \"env\""},
	}
	var buf bytes.Buffer
	FprintLabels(&buf, v)
	if !bytes.Contains(buf.Bytes(), []byte("r1")) {
		t.Error("expected resource ID in output")
	}
}
