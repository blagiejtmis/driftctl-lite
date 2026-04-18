package drift

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWatch_MaxRuns(t *testing.T) {
	state := `{"version":4,"terraform_version":"1.5.0","resources":[]}`
	f := writeTempState(t, state)

	var buf bytes.Buffer
	done := make(chan struct{})

	opts := WatchOptions{
		Interval:  10 * time.Millisecond,
		MaxRuns:   2,
		StateFile: f,
	}

	err := Watch(opts, &buf, done)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Run #1") {
		t.Error("expected Run #1 in output")
	}
	if !strings.Contains(output, "Run #2") {
		t.Error("expected Run #2 in output")
	}
}

func TestWatch_DoneChannel(t *testing.T) {
	state := `{"version":4,"terraform_version":"1.5.0","resources":[]}`
	f := writeTempState(t, state)

	var buf bytes.Buffer
	done := make(chan struct{})

	opts := WatchOptions{
		Interval:  50 * time.Millisecond,
		MaxRuns:   0,
		StateFile: f,
	}

	go func() {
		time.Sleep(80 * time.Millisecond)
		close(done)
	}()

	err := Watch(opts, &buf, done)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatch_BadStateFile(t *testing.T) {
	var buf bytes.Buffer
	done := make(chan struct{})

	opts := WatchOptions{
		Interval:  10 * time.Millisecond,
		MaxRuns:   1,
		StateFile: "/nonexistent/path/state.tfstate",
	}

	err := Watch(opts, &buf, done)
	if err == nil {
		t.Fatal("expected error for missing state file")
	}
}
