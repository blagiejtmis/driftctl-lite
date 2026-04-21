package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseRemediationReport() Report {
	return Report{
		Managed:   []Resource{{Type: "aws_s3_bucket", ID: "managed-bucket"}},
		Missing:   []Resource{{Type: "aws_instance", ID: "i-missing"}},
		Untracked: []Resource{{Type: "aws_lambda_function", ID: "fn-untracked"}},
	}
}

func TestRemediate_NoDrift(t *testing.T) {
	report := Report{
		Managed: []Resource{{Type: "aws_s3_bucket", ID: "b1"}},
	}
	hints := Remediate(report)
	if len(hints) != 0 {
		t.Errorf("expected 0 hints, got %d", len(hints))
	}
}

func TestRemediate_MissingAndUntracked(t *testing.T) {
	report := baseRemediationReport()
	hints := Remediate(report)
	if len(hints) != 2 {
		t.Fatalf("expected 2 hints, got %d", len(hints))
	}
	if hints[0].Action != "import" {
		t.Errorf("expected first hint action 'import', got %s", hints[0].Action)
	}
	if hints[1].Action != "remove" {
		t.Errorf("expected second hint action 'remove', got %s", hints[1].Action)
	}
}

func TestRemediate_HintResourceDetails(t *testing.T) {
	report := baseRemediationReport()
	hints := Remediate(report)
	if len(hints) < 1 {
		t.Fatal("expected at least one hint")
	}
	importHint := hints[0]
	if importHint.Resource.Type != "aws_instance" || importHint.Resource.ID != "i-missing" {
		t.Errorf("import hint has wrong resource: got type=%s id=%s", importHint.Resource.Type, importHint.Resource.ID)
	}
	removeHint := hints[1]
	if removeHint.Resource.Type != "aws_lambda_function" || removeHint.Resource.ID != "fn-untracked" {
		t.Errorf("remove hint has wrong resource: got type=%s id=%s", removeHint.Resource.Type, removeHint.Resource.ID)
	}
}

func TestFprintRemediation_NoHints(t *testing.T) {
	var buf bytes.Buffer
	FprintRemediation(&buf, nil)
	if !strings.Contains(buf.String(), "No remediation needed") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestFprintRemediation_WithHints(t *testing.T) {
	report := baseRemediationReport()
	hints := Remediate(report)
	var buf bytes.Buffer
	FprintRemediation(&buf, hints)
	out := buf.String()
	if !strings.Contains(out, "[IMPORT]") {
		t.Errorf("expected IMPORT in output, got: %s", out)
	}
	if !strings.Contains(out, "[REMOVE]") {
		t.Errorf("expected REMOVE in output, got: %s", out)
	}
	if !strings.Contains(out, "terraform import") {
		t.Errorf("expected terraform import hint in output")
	}
}
