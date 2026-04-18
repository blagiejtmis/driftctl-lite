package drift

import (
	"testing"
)

func makeRes(rtype, id string) Resource {
	return Resource{Type: rtype, ID: id}
}

func TestFilter_NoOptions(t *testing.T) {
	resources := []Resource{makeRes("aws_s3_bucket", "a"), makeRes("aws_instance", "b")}
	got := Filter(resources, FilterOptions{})
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilter_IncludeTypes(t *testing.T) {
	resources := []Resource{
		makeRes("aws_s3_bucket", "a"),
		makeRes("aws_instance", "b"),
		makeRes("aws_s3_bucket", "c"),
	}
	got := Filter(resources, FilterOptions{ResourceTypes: []string{"aws_s3_bucket"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	for _, r := range got {
		if r.Type != "aws_s3_bucket" {
			t.Errorf("unexpected type %s", r.Type)
		}
	}
}

func TestFilter_ExcludeTypes(t *testing.T) {
	resources := []Resource{
		makeRes("aws_s3_bucket", "a"),
		makeRes("aws_instance", "b"),
	}
	got := Filter(resources, FilterOptions{ExcludeTypes: []string{"aws_instance"}})
	if len(got) != 1 || got[0].Type != "aws_s3_bucket" {
		t.Fatalf("expected only aws_s3_bucket, got %+v", got)
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	resources := []Resource{makeRes("AWS_S3_Bucket", "a"), makeRes("aws_instance", "b")}
	got := Filter(resources, FilterOptions{ResourceTypes: []string{"aws_s3_bucket"}})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilter_ExcludeTakesPrecedence(t *testing.T) {
	resources := []Resource{makeRes("aws_s3_bucket", "a")}
	got := Filter(resources, FilterOptions{
		ResourceTypes: []string{"aws_s3_bucket"},
		ExcludeTypes:  []string{"aws_s3_bucket"},
	})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}
