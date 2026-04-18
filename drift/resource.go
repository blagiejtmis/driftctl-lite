package drift

import "fmt"

// ResourceType represents the type of a cloud resource.
type ResourceType string

const (
	ResourceTypeS3Bucket  ResourceType = "aws_s3_bucket"
	ResourceTypeEC2Instance ResourceType = "aws_instance"
	ResourceTypeSGRule     ResourceType = "aws_security_group"
)

// Resource represents a single IaC-defined resource.
type Resource struct {
	Type ResourceType
	ID   string
	Attrs map[string]string
}

// String returns a human-readable representation of the resource.
func (r Resource) String() string {
	return fmt.Sprintf("%s.%s", r.Type, r.ID)
}

// DriftResult holds the outcome of comparing live state vs IaC definitions.
type DriftResult struct {
	Managed  []Resource
	Missing  []Resource
	Untracked []Resource
}

// HasDrift returns true if any missing or untracked resources were detected.
func (d *DriftResult) HasDrift() bool {
	return len(d.Missing) > 0 || len(d.Untracked) > 0
}

// Summary returns a brief drift summary string.
func (d *DriftResult) Summary() string {
	return fmt.Sprintf("managed=%d missing=%d untracked=%d",
		len(d.Managed), len(d.Missing), len(d.Untracked))
}
