package drift

import (
	"bytes"
	"testing"
)

func makeDiffResource(typ, id string, attrs map[string]string) Resource {
	return Resource{Type: typ, ID: id, Attributes: attrs}
}

func TestDiff_NoResources(t *testing.T) {
	diffs := Diff(nil, nil)
	if len(diffs) != 0 {
		t.Fatalf("expected 0 diffs, got %d", len(diffs))
	}
}

func TestDiff_NoDrift(t *testing.T) {
	attrs := map[string]string{"region": "us-east-1"}
	live := []Resource{makeDiffResource("aws_s3_bucket", "my-bucket", attrs)}
	iac := []Resource{makeDiffResource("aws_s3_bucket", "my-bucket", attrs)}
	diffs := Diff(live, iac)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiff_AttributeChanged(t *testing.T) {
	live := []Resource{makeDiffResource("aws_s3_bucket", "my-bucket", map[string]string{"region": "us-west-2"})}
	iac := []Resource{makeDiffResource("aws_s3_bucket", "my-bucket", map[string]string{"region": "us-east-1"})}
	diffs := Diff(live, iac)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if len(diffs[0].Attributes) != 1 {
		t.Fatalf("expected 1 attribute diff, got %d", len(diffs[0].Attributes))
	}
	a := diffs[0].Attributes[0]
	if a.Key != "region" || a.IaC != "us-east-1" || a.Live != "us-west-2" {
		t.Errorf("unexpected attribute diff: %+v", a)
	}
}

func TestDiff_MissingResourceSkipped(t *testing.T) {
	// resource in IaC but not live — should not appear in diff (handled by Compare)
	iac := []Resource{makeDiffResource("aws_instance", "i-123", map[string]string{"ami": "ami-abc"})}
	diffs := Diff(nil, iac)
	if len(diffs) != 0 {
		t.Fatalf("expected 0 diffs for missing resource, got %d", len(diffs))
	}
}

func TestFprintDiff_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	FprintDiff(&buf, nil)
	if buf.String() != "No attribute drift detected.\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFprintDiff_WithDrift(t *testing.T) {
	diffs := []ResourceDiff{
		{
			Resource: makeDiffResource("aws_s3_bucket", "my-bucket", nil),
			Attributes: []AttributeDiff{
				{Key: "region", IaC: "us-east-1", Live: "us-west-2"},
			},
		},
	}
	var buf bytes.Buffer
	FprintDiff(&buf, diffs)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("my-bucket")) {
		t.Error("expected resource ID in output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("region")) {
		t.Error("expected attribute key in output")
	}
}
