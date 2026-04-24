package drift

import (
	"fmt"
	"io"
	"strings"
)

var classifyIcon = map[string]string{
	ClassCritical: "🔴",
	ClassHigh:     "🟠",
	ClassMedium:   "🟡",
	ClassLow:      "🟢",
}

// FprintClassify writes a human-readable classification report to w.
func FprintClassify(w io.Writer, r ClassifyReport) {
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "classify: no drifted resources to classify.")
		return
	}
	fmt.Fprintf(w, "%-10s %-30s %-20s %s\n", "SEVERITY", "ID", "TYPE", "REASON")
	fmt.Fprintln(w, strings.Repeat("-", 80))
	for _, res := range r.Results {
		icon := classifyIcon[res.Severity]
		fmt.Fprintf(w, "%-10s %-30s %-20s %s\n",
			icon+" "+res.Severity,
			res.Resource.ID,
			res.Resource.Type,
			res.Reason,
		)
	}
	fmt.Fprintf(w, "\nTotal classified: %d\n", len(r.Results))
}
