package drift

import (
	"fmt"
	"io"
	"time"
)

// FprintSchedule writes a human-readable schedule status to w.
func FprintSchedule(w io.Writer, cfg ScheduleConfig, now time.Time) {
	fmt.Fprintln(w, "=== Schedule Status ===")
	if !cfg.Enabled {
		fmt.Fprintln(w, "  Scheduling: disabled")
		return
	}
	fmt.Fprintf(w, "  Scheduling:   enabled\n")
	fmt.Fprintf(w, "  Interval:     %d minutes\n", cfg.IntervalMins)
	fmt.Fprintf(w, "  State file:   %s\n", cfg.StateFile)
	fmt.Fprintf(w, "  Output format:%s\n", cfg.OutputFormat)
	if cfg.LastRun.IsZero() {
		fmt.Fprintln(w, "  Last run:     never")
	} else {
		fmt.Fprintf(w, "  Last run:     %s\n", cfg.LastRun.Format(time.RFC3339))
	}
	next := NextRunTime(cfg.LastRun, cfg.IntervalMins)
	if cfg.LastRun.IsZero() {
		fmt.Fprintln(w, "  Next run:     now")
	} else {
		fmt.Fprintf(w, "  Next run:     %s\n", next.Format(time.RFC3339))
	}
	if IsDue(cfg, now) {
		fmt.Fprintln(w, "  Status:       DUE")
	} else {
		remaining := next.Sub(now).Round(time.Second)
		fmt.Fprintf(w, "  Status:       in %s\n", remaining)
	}
}

// ScheduleHasDue returns true if a scan is currently due.
func ScheduleHasDue(cfg ScheduleConfig, now time.Time) bool {
	return IsDue(cfg, now)
}
