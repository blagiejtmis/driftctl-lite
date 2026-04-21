package drift

import (
	"bytes"
	"strings"
	"testing"
)

var ownershipRules = []OwnershipRule{
	{Type: "aws_s3_bucket", Team: "platform", Email: "platform@example.com"},
	{Type: "aws_iam_role", Team: "security", Email: "security@example.com"},
}

func baseOwnershipReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "bucket-1", Type: "aws_s3_bucket"},
		},
		Missing: []Resource{
			{ID: "role-1", Type: "aws_iam_role"},
		},
		Untracked: []Resource{
			{ID: "lambda-1", Type: "aws_lambda_function"},
		},
	}
}

func TestAssignOwnership_NoRules(t *testing.T) {
	report := baseOwnershipReport()
	results := AssignOwnership(report, nil)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Unowned {
			t.Errorf("expected all resources to be unowned without rules")
		}
	}
}

func TestAssignOwnership_MatchedAndUnmatched(t *testing.T) {
	report := baseOwnershipReport()
	results := AssignOwnership(report, ownershipRules)

	ownedCount := 0
	unownedCount := 0
	for _, r := range results {
		if r.Unowned {
			unownedCount++
		} else {
			ownedCount++
		}
	}
	if ownedCount != 2 {
		t.Errorf("expected 2 owned, got %d", ownedCount)
	}
	if unownedCount != 1 {
		t.Errorf("expected 1 unowned, got %d", unownedCount)
	}
}

func TestOwnershipHasUnowned_True(t *testing.T) {
	report := baseOwnershipReport()
	results := AssignOwnership(report, ownershipRules)
	if !OwnershipHasUnowned(results) {
		t.Error("expected HasUnowned to be true")
	}
}

func TestOwnershipHasUnowned_False(t *testing.T) {
	report := Report{
		Managed: []Resource{{ID: "bucket-1", Type: "aws_s3_bucket"}},
	}
	results := AssignOwnership(report, ownershipRules)
	if OwnershipHasUnowned(results) {
		t.Error("expected HasUnowned to be false")
	}
}

func TestFprintOwnership_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintOwnership(&buf, nil)
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintOwnership_Output(t *testing.T) {
	report := baseOwnershipReport()
	results := AssignOwnership(report, ownershipRules)
	var buf bytes.Buffer
	FprintOwnership(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "platform") {
		t.Errorf("expected 'platform' team in output")
	}
	if !strings.Contains(out, "(unowned)") {
		t.Errorf("expected '(unowned)' in output")
	}
}
