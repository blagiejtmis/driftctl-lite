package drift

import (
	"fmt"
	"io"
	"strings"
)

// RemediationHint provides a suggested action for a drifted resource.
type RemediationHint struct {
	Resource Resource
	Action   string
	Hint     string
}

// Remediate generates remediation hints for a drift report.
func Remediate(report Report) []RemediationHint {
	var hints []RemediationHint

	for _, r := range report.Missing {
		hints = append(hints, RemediationHint{
			Resource: r,
			Action:   "import",
			Hint:     fmt.Sprintf("Run: terraform import %s.%s <id>", r.Type, r.ID),
		})
	}

	for _, r := range report.Untracked {
		hints = append(hints, RemediationHint{
			Resource: r,
			Action:   "remove",
			Hint:     fmt.Sprintf("Resource %s/%s exists in cloud but not in IaC — consider deleting or importing it.", r.Type, r.ID),
		})
	}

	return hints
}

// FprintRemediation writes remediation hints to the given writer.
func FprintRemediation(w io.Writer, hints []RemediationHint) {
	if len(hints) == 0 {
		fmt.Fprintln(w, "No remediation needed.")
		return
	}
	fmt.Fprintf(w, "Remediation hints (%d):\n", len(hints))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, h := range hints {
		fmt.Fprintf(w, "[%s] %s/%s\n", strings.ToUpper(h.Action), h.Resource.Type, h.Resource.ID)
		fmt.Fprintf(w, "  => %s\n", h.Hint)
	}
}
