package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseAlertReport() Report {
	return Report{
		Managed: []Resource{
			{ID: "vpc-1", Type: "aws_vpc"},
		},
		Missing: []Resource{
			{ID: "sg-1", Type: "aws_security_group"},
		},
		Untracked: []Resource{
			{ID: "s3-1", Type: "aws_s3_bucket"},
		},
	}
}

func TestEvaluateAlerts_NoResources(t *testing.T) {
	r := Report{}
	result := EvaluateAlerts(r, 5)
	if len(result.Alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(result.Alerts))
	}
}

func TestEvaluateAlerts_WarningForMissing(t *testing.T) {
	r := baseAlertReport()
	result := EvaluateAlerts(r, 5)
	if result.TotalWarn != 1 {
		t.Errorf("expected 1 warning, got %d", result.TotalWarn)
	}
	if result.TotalInfo != 1 {
		t.Errorf("expected 1 info, got %d", result.TotalInfo)
	}
}

func TestEvaluateAlerts_CriticalThreshold(t *testing.T) {
	r := Report{
		Missing: []Resource{
			{ID: "sg-1", Type: "aws_security_group"},
			{ID: "sg-2", Type: "aws_security_group"},
		},
	}
	result := EvaluateAlerts(r, 2)
	if result.TotalCrit != 2 {
		t.Errorf("expected 2 critical, got %d", result.TotalCrit)
	}
}

func TestAlertHasCritical_True(t *testing.T) {
	r := Report{
		Missing: []Resource{{ID: "sg-1", Type: "aws_security_group"}, {ID: "sg-2", Type: "aws_security_group"}},
	}
	result := EvaluateAlerts(r, 2)
	if !AlertHasCritical(result) {
		t.Error("expected HasCritical to be true")
	}
}

func TestAlertHasCritical_False(t *testing.T) {
	r := baseAlertReport()
	result := EvaluateAlerts(r, 10)
	if AlertHasCritical(result) {
		t.Error("expected HasCritical to be false")
	}
}

func TestFprintAlerts_NoAlerts(t *testing.T) {
	var buf bytes.Buffer
	FprintAlerts(&buf, AlertResult{})
	if !strings.Contains(buf.String(), "No alerts") {
		t.Errorf("expected 'No alerts', got: %s", buf.String())
	}
}

func TestFprintAlerts_WithAlerts(t *testing.T) {
	r := baseAlertReport()
	result := EvaluateAlerts(r, 10)
	var buf bytes.Buffer
	FprintAlerts(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "warning") {
		t.Errorf("expected 'warning' in output, got: %s", out)
	}
	if !strings.Contains(out, "sg-1") {
		t.Errorf("expected resource id in output, got: %s", out)
	}
}
