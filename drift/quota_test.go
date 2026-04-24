package drift

import (
	"bytes"
	"testing"
)

func baseQuotaLimits() []QuotaLimit {
	return []QuotaLimit{
		{Type: "aws_instance", SoftLimit: 3, HardLimit: 5},
		{Type: "aws_s3_bucket", SoftLimit: 2, HardLimit: 4},
	}
}

func makeQuotaResources(typ string, n int) []Resource {
	var out []Resource
	for i := 0; i < n; i++ {
		out = append(out, Resource{ID: fmt.Sprintf("%s-%d", typ, i), Type: typ})
	}
	return out
}

func TestEvaluateQuota_NoResources(t *testing.T) {
	r := EvaluateQuota(nil, baseQuotaLimits())
	if len(r.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(r.Results))
	}
}

func TestEvaluateQuota_BelowLimits(t *testing.T) {
	res := makeQuotaResources("aws_instance", 2)
	r := EvaluateQuota(res, baseQuotaLimits())
	if len(r.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(r.Results))
	}
	if r.Results[0].Exceeded != "none" {
		t.Errorf("expected none, got %s", r.Results[0].Exceeded)
	}
}

func TestEvaluateQuota_SoftExceeded(t *testing.T) {
	res := makeQuotaResources("aws_instance", 3)
	r := EvaluateQuota(res, baseQuotaLimits())
	if r.Results[0].Exceeded != "soft" {
		t.Errorf("expected soft, got %s", r.Results[0].Exceeded)
	}
}

func TestEvaluateQuota_HardExceeded(t *testing.T) {
	res := makeQuotaResources("aws_instance", 5)
	r := EvaluateQuota(res, baseQuotaLimits())
	if r.Results[0].Exceeded != "hard" {
		t.Errorf("expected hard, got %s", r.Results[0].Exceeded)
	}
}

func TestEvaluateQuota_UnknownTypeIgnored(t *testing.T) {
	res := makeQuotaResources("aws_lambda_function", 10)
	r := EvaluateQuota(res, baseQuotaLimits())
	if len(r.Results) != 0 {
		t.Errorf("expected 0 results for unknown type, got %d", len(r.Results))
	}
}

func TestQuotaHasViolations_True(t *testing.T) {
	r := QuotaReport{Results: []QuotaResult{{Exceeded: "soft"}}}
	if !QuotaHasViolations(r) {
		t.Error("expected violations")
	}
}

func TestQuotaHasViolations_False(t *testing.T) {
	r := QuotaReport{Results: []QuotaResult{{Exceeded: "none"}}}
	if QuotaHasViolations(r) {
		t.Error("expected no violations")
	}
}

func TestFprintQuota_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintQuota(&buf, QuotaReport{})
	if !bytes.Contains(buf.Bytes(), []byte("no tracked")) {
		t.Error("expected empty message")
	}
}

func TestFprintQuota_WithResults(t *testing.T) {
	var buf bytes.Buffer
	r := QuotaReport{Results: []QuotaResult{
		{Type: "aws_instance", Count: 6, SoftLimit: 3, HardLimit: 5, Exceeded: "hard"},
	}}
	FprintQuota(&buf, r)
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("HARD LIMIT")) {
		t.Errorf("expected HARD LIMIT in output, got: %s", out)
	}
}
