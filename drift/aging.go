package drift

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// AgingEntry holds resource aging information.
type AgingEntry struct {
	ResourceID string
	Type       string
	AgeDays    int
	Status     string // "managed", "missing", "untracked"
}

// AgingReport holds the full aging analysis.
type AgingReport struct {
	Entries    []AgingEntry
	AsOf       time.Time
}

// EvaluateAging computes how long each resource has been in its current drift state
// using the audit log as a source of first-seen timestamps.
func EvaluateAging(report Report, audit []AuditEntry) AgingReport {
	firstSeen := map[string]time.Time{}
	for _, e := range audit {
		for _, r := range e.Missing {
			key := r.Type + "/" + r.ID
			if _, ok := firstSeen[key]; !ok {
				firstSeen[key] = e.Timestamp
			}
		}
		for _, r := range e.Untracked {
			key := r.Type + "/" + r.ID
			if _, ok := firstSeen[key]; !ok {
				firstSeen[key] = e.Timestamp
			}
		}
	}

	now := time.Now().UTC()
	var entries []AgingEntry

	for _, r := range report.Missing {
		key := r.Type + "/" + r.ID
		age := 0
		if t, ok := firstSeen[key]; ok {
			age = int(now.Sub(t).Hours() / 24)
		}
		entries = append(entries, AgingEntry{ResourceID: r.ID, Type: r.Type, AgeDays: age, Status: "missing"})
	}
	for _, r := range report.Untracked {
		key := r.Type + "/" + r.ID
		age := 0
		if t, ok := firstSeen[key]; ok {
			age = int(now.Sub(t).Hours() / 24)
		}
		entries = append(entries, AgingEntry{ResourceID: r.ID, Type: r.Type, AgeDays: age, Status: "untracked"})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].AgeDays > entries[j].AgeDays
	})

	return AgingReport{Entries: entries, AsOf: now}
}

// AgingHasEntries returns true when there is at least one aging entry.
func AgingHasEntries(r AgingReport) bool {
	return len(r.Entries) > 0
}

// FprintAging writes a human-readable aging report to w.
func FprintAging(w io.Writer, r AgingReport) {
	if !AgingHasEntries(r) {
		fmt.Fprintln(w, "[aging] no drifted resources found")
		return
	}
	fmt.Fprintf(w, "[aging] %d drifted resource(s) as of %s\n", len(r.Entries), r.AsOf.Format("2006-01-02"))
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-12s %-30s status=%-10s age=%d days\n", e.Type, e.ResourceID, e.Status, e.AgeDays)
	}
}
