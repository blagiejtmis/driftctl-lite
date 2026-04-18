package drift

import (
	"bytes"
	"testing"
)

func baseNotifySummary(hasDrift bool) Summary {
	if hasDrift {
		return Summary{HasDrift: true, Missing: 2, Untracked: 1, Managed: 5}
	}
	return Summary{HasDrift: false, Managed: 5}
}

func TestNotify_SkippedWhenNoDriftAndOnDriftOnly(t *testing.T) {
	cfg := NotifyConfig{OnDriftOnly: true, WebhookURL: "http://example.com"}
	var buf bytes.Buffer
	res := Notify(cfg, baseNotifySummary(false), &buf)
	if !res.Skipped {
		t.Errorf("expected skipped, got sent=%v err=%v", res.Sent, res.Error)
	}
}

func TestNotify_SkippedWhenNoConfig(t *testing.T) {
	cfg := NotifyConfig{}
	var buf bytes.Buffer
	res := Notify(cfg, baseNotifySummary(true), &buf)
	if !res.Skipped {
		t.Errorf("expected skipped when no config")
	}
}

func TestNotify_SkippedSlackMissingChannel(t *testing.T) {
	cfg := NotifyConfig{SlackToken: "xoxb-token", Channel: ""}
	var buf bytes.Buffer
	res := Notify(cfg, baseNotifySummary(true), &buf)
	if !res.Skipped {
		t.Errorf("expected skipped when channel missing")
	}
}

func TestBuildMessage_NoDrift(t *testing.T) {
	s := baseNotifySummary(false)
	msg := buildMessage(s)
	if msg == "" {
		t.Error("expected non-empty message")
	}
	if contains := "No drift"; !containsStr(msg, contains) {
		t.Errorf("expected %q in message %q", contains, msg)
	}
}

func TestBuildMessage_WithDrift(t *testing.T) {
	s := baseNotifySummary(true)
	msg := buildMessage(s)
	if !containsStr(msg, "Missing: 2") {
		t.Errorf("expected missing count in message: %q", msg)
	}
	if !containsStr(msg, "Untracked: 1") {
		t.Errorf("expected untracked count in message: %q", msg)
	}
}
