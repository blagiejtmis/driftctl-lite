package drift

import (
	"fmt"
	"io"
)

// FprintTagViolations writes a human-readable tag compliance report to w.
func FprintTagViolations(w io.Writer, violations []TagViolation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "[tags] All resources are tag-compliant.")
		return
	}
	fmt.Fprintf(w, "[tags] %d tag violation(s) found:\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(w, "  - %s/%s: %s\n", v.Resource.Type, v.Resource.ID, v.Reason)
	}
}

// TagHasViolations returns true when the slice is non-empty.
func TagHasViolations(violations []TagViolation) bool {
	return len(violations) > 0
}
