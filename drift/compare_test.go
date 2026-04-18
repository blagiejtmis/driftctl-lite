package drift

import (
	"testing"
)

func makeResource(t ResourceType, id string) Resource {
	return Resource{Type: t, ID: id}
}

func TestCompare_NoResources(t *testing.T) {
	result := Compare(nil, nil)
	if result.HasDrift() {
		t.Error("expected no drift for empty inputs")
	}
}

func TestCompare_AllManaged(t *testing.T) {
	iac := []Resource{makeResource(ResourceTypeS3Bucket, "my-bucket")}
	live := []Resource{makeResource(ResourceTypeS3Bucket, "my-bucket")}

	result := Compare(iac, live)
	if len(result.Managed) != 1 {
		t.Errorf("expected 1 managed, got %d", len(result.Managed))
	}
	if result.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestCompare_MissingResource(t *testing.T) {
	iac := []Resource{
		makeResource(ResourceTypeS3Bucket, "my-bucket"),
		makeResource(ResourceTypeEC2Instance, "i-1234"),
	}
	live := []Resource{makeResource(ResourceTypeS3Bucket, "my-bucket")}

	result := Compare(iac, live)
	if len(result.Missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(result.Missing))
	}
	if result.Missing[0].ID != "i-1234" {
		t.Errorf("unexpected missing resource: %s", result.Missing[0].ID)
	}
}

func TestCompare_UntrackedResource(t *testing.T) {
	iac := []Resource{makeResource(ResourceTypeS3Bucket, "my-bucket")}
	live := []Resource{
		makeResource(ResourceTypeS3Bucket, "my-bucket"),
		makeResource(ResourceTypeSGRule, "sg-9999"),
	}

	result := Compare(iac, live)
	if len(result.Untracked) != 1 {
		t.Errorf("expected 1 untracked, got %d", len(result.Untracked))
	}
	if result.Untracked[0].ID != "sg-9999" {
		t.Errorf("unexpected untracked resource: %s", result.Untracked[0].ID)
	}
}

func TestDriftResult_Summary(t *testing.T) {
	d := DriftResult{
		Managed:   []Resource{makeResource(ResourceTypeS3Bucket, "b1")},
		Missing:   []Resource{makeResource(ResourceTypeEC2Instance, "i1")},
		Untracked: []Resource{},
	}
	expected := "managed=1 missing=1 untracked=0"
	if s := d.Summary(); s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
