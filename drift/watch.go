package drift

import (
	"fmt"
	"io"
	"time"
)

// WatchOptions configures the watch loop.
type WatchOptions struct {
	Interval  time.Duration
	MaxRuns   int // 0 = run forever
	StateFile string
	Filter    FilterOptions
}

// WatchResult holds the outcome of a single watch tick.
type WatchResult struct {
	Run       int
	Timestamp time.Time
	Report    Report
	Summary   Summary
}

// Watch repeatedly scans for drift at the given interval, writing results to w.
// It stops after MaxRuns iterations (if > 0) or when done is closed.
func Watch(opts WatchOptions, w io.Writer, done <-chan struct{}) error {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	run := 0
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	execute := func() error {
		run++
		resources, err := Scan(opts.StateFile)
		if err != nil {
			return fmt.Errorf("scan error on run %d: %w", run, err)
		}
		filtered := Filter(resources, opts.Filter)
		report := Compare(filtered, filtered)
		summary := Summarize(report)
		result := WatchResult{
			Run:       run,
			Timestamp: time.Now().UTC(),
			Report:    report,
			Summary:   summary,
		}
		fmt.Fprintf(w, "[%s] Run #%d\n", result.Timestamp.Format(time.RFC3339), result.Run)
		FprintSummary(w, summary)
		if opts.MaxRuns > 0 && run >= opts.MaxRuns {
			return io.EOF
		}
		return nil
	}

	// Run immediately before waiting for first tick.
	if err := execute(); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			if err := execute(); err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}
		}
	}
}
