package drift

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// DriftAgeEntry records how long a resource has been in a drifted state.
type DriftAgeEntry struct {
	ResourceID string        `json:"resource_id"`
	Type       string        `json:"type"`
	DriftKind  string        `json:"drift_kind"` // "missing" | "untracked" | "changed"
	FirstSeen  time.Time     `json:"first_seen"`
	AgeDays    float64       `json:"age_days"`
}

// DriftAgeReport holds all entries for the current evaluation.
type DriftAgeReport struct {
	Entries []DriftAgeEntry `json:"entries"`
}

// EvaluateDriftAge computes how long each drifted resource has been drifted,
// using the audit log to determine the first time the resource appeared as drifted.
func EvaluateDriftAge(report Report, audit []AuditEntry) DriftAgeReport {
	// Build a map: resourceKey -> earliest audit timestamp where it was drifted.
	firstSeen := map[string]time.Time{}
	for _, entry := range audit {
		for _, r := range entry.Missing {
			k := resourceKey(r)
			if t, ok := firstSeen[k]; !ok || entry.Timestamp.Before(t) {
				firstSeen[k] = entry.Timestamp
			}
		}
		for _, r := range entry.Untracked {
			k := resourceKey(r)
			if t, ok := firstSeen[k]; !ok || entry.Timestamp.Before(t) {
				firstSeen[k] = entry.Timestamp
			}
		}
	}

	now := time.Now().UTC()
	var entries []DriftAgeEntry

	for _, r := range report.Missing {
		fs := firstSeenOrNow(firstSeen, r, now)
		entries = append(entries, DriftAgeEntry{
			ResourceID: r.ID,
			Type:       r.Type,
			DriftKind:  "missing",
			FirstSeen:  fs,
			AgeDays:    now.Sub(fs).Hours() / 24,
		})
	}
	for _, r := range report.Untracked {
		fs := firstSeenOrNow(firstSeen, r, now)
		entries = append(entries, DriftAgeEntry{
			ResourceID: r.ID,
			Type:       r.Type,
			DriftKind:  "untracked",
			FirstSeen:  fs,
			AgeDays:    now.Sub(fs).Hours() / 24,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].AgeDays > entries[j].AgeDays
	})

	return DriftAgeReport{Entries: entries}
}

func firstSeenOrNow(m map[string]time.Time, r Resource, now time.Time) time.Time {
	if t, ok := m[resourceKey(r)]; ok {
		return t
	}
	return now
}

// DriftAgeHasEntries returns true when the report contains any entries.
func DriftAgeHasEntries(r DriftAgeReport) bool {
	return len(r.Entries) > 0
}

// FprintDriftAge writes a human-readable drift-age table to w.
func FprintDriftAge(w io.Writer, r DriftAgeReport) {
	if !DriftAgeHasEntries(r) {
		fmt.Fprintln(w, "No drift-age entries found.")
		return
	}
	fmt.Fprintf(w, "%-36s  %-20s  %-10s  %s\n", "Resource ID", "Type", "Kind", "Age (days)")
	fmt.Fprintf(w, "%s\n", "---------------------------------------------------------------------------------------------")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "%-36s  %-20s  %-10s  %.1f\n", e.ResourceID, e.Type, e.DriftKind, e.AgeDays)
	}
}
