package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseGroupReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc", Attributes: map[string]string{}},
			{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{}},
		},
		Missing: []Resource{
			{ID: "vpc-2", Type: "aws_vpc", Attributes: map[string]string{}},
		},
		Untracked: []Resource{
			{ID: "sg-99", Type: "aws_security_group", Attributes: map[string]string{}},
		},
	}
}

func TestGroupByType_Empty(t *testing.T) {
	g := GroupByType(Report{})
	if len(g.Groups) != 0 {
		t.Fatalf("expected empty groups, got %d", len(g.Groups))
	}
}

func TestGroupByType_CorrectCounts(t *testing.T) {
	g := GroupByType(baseGroupReport())

	vpc, ok := g.Groups["aws_vpc"]
	if !ok {
		t.Fatal("expected aws_vpc group")
	}
	if len(vpc.Managed) != 1 {
		t.Errorf("aws_vpc managed: want 1, got %d", len(vpc.Managed))
	}
	if len(vpc.Missing) != 1 {
		t.Errorf("aws_vpc missing: want 1, got %d", len(vpc.Missing))
	}

	sg, ok := g.Groups["aws_security_group"]
	if !ok {
		t.Fatal("expected aws_security_group group")
	}
	if len(sg.Untracked) != 1 {
		t.Errorf("aws_security_group untracked: want 1, got %d", len(sg.Untracked))
	}
}

func TestFprintGroup_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintGroup(&buf, GroupedReport{})
	if !strings.Contains(buf.String(), "No resources") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintGroup_OutputContainsTypes(t *testing.T) {
	var buf bytes.Buffer
	FprintGroup(&buf, GroupByType(baseGroupReport()))
	out := buf.String()

	if !strings.Contains(out, "AWS_VPC") {
		t.Errorf("expected AWS_VPC in output")
	}
	if !strings.Contains(out, "AWS_SECURITY_GROUP") {
		t.Errorf("expected AWS_SECURITY_GROUP in output")
	}
	if !strings.Contains(out, "vpc-2") {
		t.Errorf("expected missing resource vpc-2 in output")
	}
	if !strings.Contains(out, "sg-99") {
		t.Errorf("expected untracked resource sg-99 in output")
	}
}
