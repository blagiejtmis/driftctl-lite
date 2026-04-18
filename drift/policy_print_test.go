package drift

import (
	"bytes"
	"strings"
	"testing"
)

func basePolicyResults() []PolicyResult {
	return []PolicyResult{
		{
			Rule:       PolicyRule{ID: "p1", Severity: "error", Message: "bucket missing"},
			Violations: []Resource{{Type: "aws_s3_bucket", ID: "logs"}},
		},
		{
			Rule:       PolicyRule{ID: "p2", Severity: "warn", Message: "sg untracked"},
			Violations: []Resource{{Type: "aws_security_group", ID: "sg-1"}},
		},
	}
}

func TestFprintPolicy_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	FprintPolicy(&buf, nil)
	if !strings.Contains(buf.String(), "No policy violations") {
		t.Errorf("expected no-violation message, got: %s", buf.String())
	}
}

func TestFprintPolicy_WithResults(t *testing.T) {
	var buf bytes.Buffer
	FprintPolicy(&buf, basePolicyResults())
	out := buf.String()
	if !strings.Contains(out, "1 error(s)") {
		t.Errorf("expected error count, got: %s", out)
	}
	if !strings.Contains(out, "1 warning(s)") {
		t.Errorf("expected warn count, got: %s", out)
	}
	if !strings.Contains(out, "aws_s3_bucket.logs") {
		t.Errorf("expected violation resource, got: %s", out)
	}
}

func TestPolicyHasErrors_True(t *testing.T) {
	if !PolicyHasErrors(basePolicyResults()) {
		t.Fatal("expected errors")
	}
}

func TestPolicyHasErrors_False(t *testing.T) {
	results := []PolicyResult{
		{Rule: PolicyRule{Severity: "warn"}},
	}
	if PolicyHasErrors(results) {
		t.Fatal("expected no errors")
	}
}
