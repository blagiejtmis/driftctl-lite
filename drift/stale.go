package drift

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// StaleResource represents a resource that has not been updated recently.
type StaleResource struct {
	Resource Resource
	LastSeen time.Time
	AgeDays  int
}

// StaleReport holds the result of a staleness evaluation.
type StaleReport struct {
	Stale     []StaleResource
	Threshold int // days
}

// EvaluateStale checks which resources have a "last_modified" attribute older
// than thresholdDays days relative to now.
func EvaluateStale(resources []Resource, thresholdDays int, now time.Time) StaleReport {
	report := StaleReport{Threshold: thresholdDays}

	for _, r := range resources {
		raw, ok := r.Attributes["last_modified"]
		if !ok {
			continue
		}
		s, ok := raw.(string)
		if !ok {
			continue
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			continue
		}
		age := int(now.Sub(t).Hours() / 24)
		if age >= thresholdDays {
			report.Stale = append(report.Stale, StaleResource{
				Resource: r,
				LastSeen: t,
				AgeDays:  age,
			})
		}
	}

	sort.Slice(report.Stale, func(i, j int) bool {
		return report.Stale[i].AgeDays > report.Stale[j].AgeDays
	})
	return report
}

// StaleHasEntries returns true when there is at least one stale resource.
func StaleHasEntries(r StaleReport) bool {
	return len(r.Stale) > 0
}

// FprintStale writes a human-readable staleness report to w.
func FprintStale(w io.Writer, r StaleReport) {
	if !StaleHasEntries(r) {
		fmt.Fprintln(w, "No stale resources found.")
		return
	}
	fmt.Fprintf(w, "Stale resources (threshold: %d days):\n", r.Threshold)
	for _, s := range r.Stale {
		fmt.Fprintf(w, "  [%s] %s — %d days old (last_modified: %s)\n",
			s.Resource.Type, s.Resource.ID, s.AgeDays,
			s.LastSeen.Format(time.RFC3339))
	}
}
