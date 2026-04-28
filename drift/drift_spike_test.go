package drift

import (
	"bytes"
	"strings"
	"testing"
)

func baseSpikeData() (map[string]int, map[string]int) {
	prev := map[string]int{
		"aws_s3_bucket":       4,
		"aws_security_group":  2,
		"aws_instance":        1,
	}
	curr := map[string]int{
		"aws_s3_bucket":       4,  // no change
		"aws_security_group":  6,  // +4 = 200% spike
		"aws_instance":        1,  // no change
		"aws_iam_role":        3,  // new type
	}
	return prev, curr
}

func TestEvaluateSpike_NoResources(t *testing.T) {
	r := EvaluateSpike(nil, nil, 50.0)
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestEvaluateSpike_DetectsSpike(t *testing.T) {
	prev, curr := baseSpikeData()
	r := EvaluateSpike(prev, curr, 50.0)

	var sg *SpikeEntry
	for i := range r.Entries {
		if r.Entries[i].Type == "aws_security_group" {
			sg = &r.Entries[i]
		}
	}
	if sg == nil {
		t.Fatal("expected aws_security_group entry")
	}
	if !sg.IsSpike {
		t.Errorf("expected aws_security_group to be flagged as spike")
	}
	if sg.Delta != 4 {
		t.Errorf("expected delta 4, got %d", sg.Delta)
	}
}

func TestEvaluateSpike_NewTypeIsSpike(t *testing.T) {
	prev, curr := baseSpikeData()
	r := EvaluateSpike(prev, curr, 50.0)

	var iam *SpikeEntry
	for i := range r.Entries {
		if r.Entries[i].Type == "aws_iam_role" {
			iam = &r.Entries[i]
		}
	}
	if iam == nil {
		t.Fatal("expected aws_iam_role entry")
	}
	if !iam.IsSpike {
		t.Errorf("expected aws_iam_role (new type) to be flagged as spike")
	}
}

func TestEvaluateSpike_NoSpikeWhenBelowThreshold(t *testing.T) {
	prev := map[string]int{"aws_s3_bucket": 10}
	curr := map[string]int{"aws_s3_bucket": 11}
	r := EvaluateSpike(prev, curr, 50.0)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if r.Entries[0].IsSpike {
		t.Errorf("expected no spike for 10%% change below 50%% threshold")
	}
}

func TestSpikeHasSpikes_True(t *testing.T) {
	r := SpikeReport{Entries: []SpikeEntry{{IsSpike: true}}}
	if !SpikeHasSpikes(r) {
		t.Error("expected SpikeHasSpikes to return true")
	}
}

func TestSpikeHasSpikes_False(t *testing.T) {
	r := SpikeReport{Entries: []SpikeEntry{{IsSpike: false}}}
	if SpikeHasSpikes(r) {
		t.Error("expected SpikeHasSpikes to return false")
	}
}

func TestFprintSpike_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintSpike(&buf, SpikeReport{})
	if !strings.Contains(buf.String(), "No drift spike") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFprintSpike_OutputContainsTypes(t *testing.T) {
	prev, curr := baseSpikeData()
	r := EvaluateSpike(prev, curr, 50.0)
	var buf bytes.Buffer
	FprintSpike(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "aws_security_group") {
		t.Errorf("expected aws_security_group in output")
	}
	if !strings.Contains(out, "SPIKE") {
		t.Errorf("expected SPIKE label in output")
	}
}
