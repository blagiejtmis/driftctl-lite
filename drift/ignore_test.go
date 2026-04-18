package drift

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "driftignore-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadIgnoreFile_Missing(t *testing.T) {
	rules, err := LoadIgnoreFile("/nonexistent/.driftignore")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected no rules, got %d", len(rules))
	}
}

func TestLoadIgnoreFile_ParsesRules(t *testing.T) {
	path := writeTempIgnore(t, "# comment\naws_s3_bucket\naws_instance/i-1234\n\naws_vpc/vpc-abc\n")
	rules, err := LoadIgnoreFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Type != "aws_s3_bucket" || rules[0].ID != "" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Type != "aws_instance" || rules[1].ID != "i-1234" {
		t.Errorf("unexpected rule[1]: %+v", rules[1])
	}
}

func TestApplyIgnore_NoRules(t *testing.T) {
	res := []Resource{{Type: "aws_s3_bucket", ID: "my-bucket"}}
	out := ApplyIgnore(res, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(out))
	}
}

func TestApplyIgnore_ByType(t *testing.T) {
	res := []Resource{
		{Type: "aws_s3_bucket", ID: "b1"},
		{Type: "aws_instance", ID: "i-1"},
	}
	rules := []IgnoreRule{{Type: "aws_s3_bucket"}}
	out := ApplyIgnore(res, rules)
	if len(out) != 1 || out[0].Type != "aws_instance" {
		t.Errorf("expected only aws_instance, got %+v", out)
	}
}

func TestApplyIgnore_ByTypeAndID(t *testing.T) {
	res := []Resource{
		{Type: "aws_instance", ID: "i-1"},
		{Type: "aws_instance", ID: "i-2"},
	}
	rules := []IgnoreRule{{Type: "aws_instance", ID: "i-1"}}
	out := ApplyIgnore(res, rules)
	if len(out) != 1 || out[0].ID != "i-2" {
		t.Errorf("expected only i-2, got %+v", out)
	}
}

func TestApplyIgnore_CaseInsensitive(t *testing.T) {
	res := []Resource{{Type: "AWS_S3_Bucket", ID: "b1"}}
	rules := []IgnoreRule{{Type: "aws_s3_bucket"}}
	out := ApplyIgnore(res, rules)
	if len(out) != 0 {
		t.Errorf("expected resource to be ignored, got %+v", out)
	}
}
