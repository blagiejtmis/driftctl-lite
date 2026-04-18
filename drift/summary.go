package drift

import (
	"fmt"
	"io"
)

// Summary holds aggregated drift statistics for a scan report.
type Summary struct {
	Total     int
	Managed   int
	Missing   int
	Untracked int
	Changed   int
}

// Summarize computes a Summary from a Report.
func Summarize(r Report) Summary {
	s := Summary{}
	s.Missing = len(r.Missing)
	s.Untracked = len(r.Untracked)
	s.Changed = len(r.Changed)
	s.Managed = len(r.Managed)
	s.Total = s.Managed + s.Missing + s.Untracked + s.Changed
	return s
}

// HasDrift returns true when any drift is detected.
func (s Summary) HasDrift() bool {
	return s.Missing > 0 || s.Untracked > 0 || s.Changed > 0
}

// FprintSummary writes a human-readable summary to w.
func FprintSummary(w io.Writer, s Summary) {
	fmt.Fprintf(w, "Drift Summary\n")
	fmt.Fprintf(w, "  Total resources : %d\n", s.Total)
	fmt.Fprintf(w, "  Managed         : %d\n", s.Managed)
	fmt.Fprintf(w, "  Missing         : %d\n", s.Missing)
	fmt.Fprintf(w, "  Untracked       : %d\n", s.Untracked)
	fmt.Fprintf(w, "  Changed         : %d\n", s.Changed)
	if s.HasDrift() {
		fmt.Fprintf(w, "\nDrift detected!\n")
	} else {
		fmt.Fprintf(w, "\nNo drift detected.\n")
	}
}
