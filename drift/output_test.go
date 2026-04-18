package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeReport() Report {
	return Report{
		Managed:   []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}},
		Missing:   []Resource{{Type: "aws_instance", ID: "i-123"}},
		Untracked: []Resource{{Type: "aws_s3_bucket", ID: "other-bucket"}},
	}
}

func TestReport_HasDrift(t *testing.T) {
	r := makeReport()
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}

	clean := Report{Managed: r.Managed}
	if clean.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestReport_FprintText(t *testing.T) {
	var buf bytes.Buffer
	r := makeReport()
	if err := r.Fprint(&buf, OutputText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING in output")
	}
	if !strings.Contains(out, "UNTRACKED") {
		t.Error("expected UNTRACKED in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK in output")
	}
	if !strings.Contains(out, "i-123") {
		t.Error("expected resource ID in output")
	}
}

func TestReport_FprintJSON(t *testing.T) {
	var buf bytes.Buffer
	r := makeReport()
	if err := r.Fprint(&buf, OutputJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded.Missing) != 1 || decoded.Missing[0].ID != "i-123" {
		t.Errorf("unexpected missing resources: %+v", decoded.Missing)
	}
}
