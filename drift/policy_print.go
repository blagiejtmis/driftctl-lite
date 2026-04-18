package drift

import (
	"fmt"
	"io"
	"strings"
)

// FprintPolicy writes policy evaluation results in human-readable form.
func FprintPolicy(w io.Writer, results []PolicyResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "✅  No policy violations found.")
		return
	}

	errorCount := 0
	warnCount := 0
	for _, r := range results {
		switch strings.ToLower(r.Rule.Severity) {
		case "error":
			errorCount++
		case "warn":
			warnCount++
		}
	}

	fmt.Fprintf(w, "Policy Results: %d error(s), %d warning(s)\n", errorCount, warnCount)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, r := range results {
		sev := strings.ToUpper(r.Rule.Severity)
		fmt.Fprintf(w, "[%s] %s: %s\n", sev, r.Rule.ID, r.Rule.Message)
		for _, v := range r.Violations {
			fmt.Fprintf(w, "  - %s.%s\n", v.Type, v.ID)
		}
	}
}

// PolicyHasErrors returns true if any result has severity "error".
func PolicyHasErrors(results []PolicyResult) bool {
	for _, r := range results {
		if strings.EqualFold(r.Rule.Severity, "error") {
			return true
		}
	}
	return false
}
