package drift

import (
	"bytes"
	"testing"
)

func baseDepsResources() []Resource {
	return []Resource{
		{
			ID:   "sg-1",
			Type: "aws_security_group",
			Attributes: map[string]interface{}{
				"name": "web-sg",
			},
		},
		{
			ID:   "inst-1",
			Type: "aws_instance",
			Attributes: map[string]interface{}{
				"depends_on": "aws_security_group.sg-1",
			},
		},
		{
			ID:   "lb-1",
			Type: "aws_lb",
			Attributes: map[string]interface{}{
				"depends_on": "aws_instance.inst-1, aws_security_group.sg-1",
			},
		},
	}
}

func TestBuildDepGraph_NoDepends(t *testing.T) {
	resources := []Resource{{ID: "r1", Type: "aws_s3_bucket", Attributes: map[string]interface{}{}}}
	g := BuildDepGraph(resources)
	if len(g.Edges) != 0 {
		t.Fatalf("expected 0 edges, got %d", len(g.Edges))
	}
}

func TestBuildDepGraph_SingleDep(t *testing.T) {
	resources := baseDepsResources()
	g := BuildDepGraph(resources)
	// inst-1 depends on sg-1, lb-1 depends on inst-1 and sg-1 → 3 edges total
	if len(g.Edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(g.Edges))
	}
}

func TestAffected_DirectAndTransitive(t *testing.T) {
	resources := baseDepsResources()
	g := BuildDepGraph(resources)
	// Changing sg-1 should affect inst-1 and lb-1 (both depend on it).
	affected := g.Affected("aws_security_group.sg-1")
	if len(affected) != 2 {
		t.Fatalf("expected 2 affected resources, got %d: %v", len(affected), affected)
	}
}

func TestAffected_NoDependents(t *testing.T) {
	resources := baseDepsResources()
	g := BuildDepGraph(resources)
	affected := g.Affected("aws_lb.lb-1")
	if len(affected) != 0 {
		t.Fatalf("expected 0 affected, got %d", len(affected))
	}
}

func TestFprintDeps_Empty(t *testing.T) {
	g := &DepGraph{adj: make(map[string][]string)}
	var buf bytes.Buffer
	FprintDeps(&buf, g)
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output for empty graph")
	}
	if !bytes.Contains(buf.Bytes(), []byte("No resource")) {
		t.Errorf("expected 'No resource' message, got: %s", buf.String())
	}
}

func TestFprintDeps_WithEdges(t *testing.T) {
	resources := baseDepsResources()
	g := BuildDepGraph(resources)
	var buf bytes.Buffer
	FprintDeps(&buf, g)
	if !bytes.Contains(buf.Bytes(), []byte("-->")) {
		t.Errorf("expected arrow in output, got: %s", buf.String())
	}
}
