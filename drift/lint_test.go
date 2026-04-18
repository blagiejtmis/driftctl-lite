package drift

import (
	"strings"
	"testing"
)

func baseLintResources() []Resource {
	return []Resource{
		{ID: "i-123", Type: "aws_instance", Attributes: map[string]string{"ami": "ami-abc"}},
	}
}

func TestLint_NoViolations(t *testing.T) {
	r := Lint(baseLintResources())
	if LintHasErrors(r) {
		t.Fatalf("expected no violations, got %d", len(r.Results))
	}
}

func TestLint_MissingID(t *testing.T) {
	res := []Resource{{ID: "", Type: "aws_instance", Attributes: map[string]string{"x": "y"}}}
	r := Lint(res)
	if !LintHasErrors(r) {
		t.Fatal("expected violations")
	}
	if r.Results[0].Violations[0].Name != "missing-id" {
		t.Errorf("unexpected rule: %s", r.Results[0].Violations[0].Name)
	}
}

func TestLint_MissingType(t *testing.T) {
	res := []Resource{{ID: "i-1", Type: "", Attributes: map[string]string{"x": "y"}}}
	r := Lint(res)
	if !LintHasErrors(r) {
		t.Fatal("expected violations")
	}
	found := false
	for _, v := range r.Results[0].Violations {
		if v.Name == "missing-type" {
			found = true
		}
	}
	if !found {
		t.Error("missing-type rule not found")
	}
}

func TestLint_NoAttributes(t *testing.T) {
	res := []Resource{{ID: "i-1", Type: "aws_s3_bucket", Attributes: map[string]string{}}}
	r := Lint(res)
	if !LintHasErrors(r) {
		t.Fatal("expected violations")
	}
	if r.Results[0].Violations[0].Name != "no-attributes" {
		t.Errorf("unexpected rule: %s", r.Results[0].Violations[0].Name)
	}
}

func TestFprintLint_NoIssues(t *testing.T) {
	r := Lint(baseLintResources())
	var sb strings.Builder
	FprintLint(&sb, r)
	if !strings.Contains(sb.String(), "no issues") {
		t.Errorf("unexpected output: %s", sb.String())
	}
}

func TestFprintLint_WithIssues(t *testing.T) {
	res := []Resource{{ID: "", Type: "", Attributes: map[string]string{}}}
	r := Lint(res)
	var sb strings.Builder
	FprintLint(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "missing-id") {
		t.Errorf("expected missing-id in output: %s", out)
	}
}
