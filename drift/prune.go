package drift

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// PruneEntry represents a resource candidate for pruning.
type PruneEntry struct {
	Resource Resource
	Reason   string
	AgeDays  int
}

// PruneReport holds the result of a prune evaluation.
type PruneReport struct {
	Entries   []PruneEntry
	Evaluated int
}

// PruneOptions controls which resources are flagged for pruning.
type PruneOptions struct {
	MaxAgeDays    int
	Types         []string
	OnlyUntracked bool
}

// EvaluatePrune scans resources and flags candidates for pruning based on age
// and drift status. Untracked resources older than MaxAgeDays are included.
func EvaluatePrune(report Report, auditPath string, opts PruneOptions) (PruneReport, error) {
	history, _ := LoadAudit(auditPath)

	ageMap := map[string]int{}
	for _, e := range history {
		for _, r := range e.Missing {
			key := resourceKey(r)
			days := int(time.Since(e.Timestamp).Hours() / 24)
			if existing, ok := ageMap[key]; !ok || days > existing {
				ageMap[key] = days
			}
		}
	}

	var entries []PruneEntry
	candidates := report.Untracked
	if !opts.OnlyUntracked {
		candidates = append(candidates, report.Missing...)
	}

	for _, r := range candidates {
		if len(opts.Types) > 0 && !containsStr(opts.Types, strings.ToLower(r.Type)) {
			continue
		}
		age := ageMap[resourceKey(r)]
		if opts.MaxAgeDays > 0 && age < opts.MaxAgeDays {
			continue
		}
		reason := "untracked"
		if age > 0 {
			reason = fmt.Sprintf("untracked for %d days", age)
		}
		entries = append(entries, PruneEntry{Resource: r, Reason: reason, AgeDays: age})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].AgeDays > entries[j].AgeDays
	})

	return PruneReport{Entries: entries, Evaluated: len(candidates)}, nil
}

// PruneHasEntries returns true when there are prune candidates.
func PruneHasEntries(pr PruneReport) bool {
	return len(pr.Entries) > 0
}

// FprintPrune writes a human-readable prune report to w.
func FprintPrune(w io.Writer, pr PruneReport) {
	if !PruneHasEntries(pr) {
		fmt.Fprintln(w, "No resources flagged for pruning.")
		return
	}
	fmt.Fprintf(w, "Prune candidates (%d of %d evaluated):\n", len(pr.Entries), pr.Evaluated)
	for _, e := range pr.Entries {
		fmt.Fprintf(w, "  [%s] %s — %s\n", e.Resource.Type, e.Resource.ID, e.Reason)
	}
}
