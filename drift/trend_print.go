package drift

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// FprintTrend writes a human-readable trend table to w.
func FprintTrend(w io.Writer, log TrendLog) {
	if len(log.Entries) == 0 {
		fmt.Fprintln(w, "No trend data recorded yet.")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tSCORE\tGRADE\tMANAGED\tMISSING\tUNTRACKED")
	for _, e := range log.Entries {
		fmt.Fprintf(tw, "%s\t%.1f%%\t%s\t%d\t%d\t%d\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Score,
			e.Grade,
			e.Managed,
			e.Missing,
			e.Untracked,
		)
	}
	tw.Flush()
}

// TrendImproving returns true if the latest score is higher than the previous.
func TrendImproving(log TrendLog) bool {
	if len(log.Entries) < 2 {
		return false
	}
	last := log.Entries[len(log.Entries)-1]
	prev := log.Entries[len(log.Entries)-2]
	return last.Score > prev.Score
}
