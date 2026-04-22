package drift

import (
	"bytes"
	"testing"
	"time"
)

func baseStaleResources(now time.Time) []Resource {
	return []Resource{
		{
			ID:   "res-1",
			Type: "aws_instance",
			Attributes: map[string]interface{}{
				"last_modified": now.Add(-100 * 24 * time.Hour).Format(time.RFC3339),
			},
		},
		{
			ID:   "res-2",
			Type: "aws_s3_bucket",
			Attributes: map[string]interface{}{
				"last_modified": now.Add(-5 * 24 * time.Hour).Format(time.RFC3339),
			},
		},
		{
			ID:         "res-3",
			Type:       "aws_lambda_function",
			Attributes: map[string]interface{}{},
		},
	}
}

func TestEvaluateStale_NoResources(t *testing.T) {
	now := time.Now()
	report := EvaluateStale([]Resource{}, 30, now)
	if StaleHasEntries(report) {
		t.Fatal("expected no stale resources")
	}
}

func TestEvaluateStale_BelowThreshold(t *testing.T) {
	now := time.Now()
	res := baseStaleResources(now)
	report := EvaluateStale(res, 200, now)
	if StaleHasEntries(report) {
		t.Fatalf("expected 0 stale, got %d", len(report.Stale))
	}
}

func TestEvaluateStale_AboveThreshold(t *testing.T) {
	now := time.Now()
	res := baseStaleResources(now)
	report := EvaluateStale(res, 30, now)
	if len(report.Stale) != 1 {
		t.Fatalf("expected 1 stale resource, got %d", len(report.Stale))
	}
	if report.Stale[0].Resource.ID != "res-1" {
		t.Errorf("expected res-1, got %s", report.Stale[0].Resource.ID)
	}
}

func TestEvaluateStale_SortedDescending(t *testing.T) {
	now := time.Now()
	res := []Resource{
		{ID: "a", Type: "t", Attributes: map[string]interface{}{"last_modified": now.Add(-40 * 24 * time.Hour).Format(time.RFC3339)}},
		{ID: "b", Type: "t", Attributes: map[string]interface{}{"last_modified": now.Add(-90 * 24 * time.Hour).Format(time.RFC3339)}},
	}
	report := EvaluateStale(res, 10, now)
	if report.Stale[0].Resource.ID != "b" {
		t.Errorf("expected b first (oldest), got %s", report.Stale[0].Resource.ID)
	}
}

func TestFprintStale_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	FprintStale(&buf, StaleReport{Threshold: 30})
	if buf.String() == "" {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("No stale")) {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestFprintStale_WithEntries(t *testing.T) {
	now := time.Now()
	res := baseStaleResources(now)
	report := EvaluateStale(res, 30, now)
	var buf bytes.Buffer
	FprintStale(&buf, report)
	if !bytes.Contains(buf.Bytes(), []byte("res-1")) {
		t.Errorf("expected res-1 in output: %s", buf.String())
	}
}
