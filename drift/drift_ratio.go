package drift

import (
	"fmt"
	"io"
	"sort"
)

// RatioEntry holds the drift ratio metrics for a single resource type.
type RatioEntry struct {
	Type       string
	Total      int
	Drifted    int
	Ratio      float64 // fraction of drifted resources (0.0 – 1.0)
	RatioPct   float64 // ratio expressed as a percentage
}

// RatioReport is the result of BuildDriftRatio.
type RatioReport struct {
	Entries      []RatioEntry
	OverallRatio float64
}

// BuildDriftRatio computes per-type and overall drift ratios from a Report.
func BuildDriftRatio(r Report) RatioReport {
	type counts struct{ total, drifted int }
	buckets := map[string]*counts{}

	for _, res := range r.Managed {
		if _, ok := buckets[res.Type]; !ok {
			buckets[res.Type] = &counts{}
		}
		buckets[res.Type].total++
	}
	for _, res := range r.Missing {
		if _, ok := buckets[res.Type]; !ok {
			buckets[res.Type] = &counts{}
		}
		buckets[res.Type].total++
		buckets[res.Type].drifted++
	}
	for _, res := range r.Untracked {
		if _, ok := buckets[res.Type]; !ok {
			buckets[res.Type] = &counts{}
		}
		buckets[res.Type].total++
		buckets[res.Type].drifted++
	}

	var entries []RatioEntry
	totalAll, driftedAll := 0, 0
	for typ, c := range buckets {
		ratio := 0.0
		if c.total > 0 {
			ratio = float64(c.drifted) / float64(c.total)
		}
		entries = append(entries, RatioEntry{
			Type:     typ,
			Total:    c.total,
			Drifted:  c.drifted,
			Ratio:    ratio,
			RatioPct: ratio * 100,
		})
		totalAll += c.total
		driftedAll += c.drifted
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Ratio != entries[j].Ratio {
			return entries[i].Ratio > entries[j].Ratio
		}
		return entries[i].Type < entries[j].Type
	})

	overall := 0.0
	if totalAll > 0 {
		overall = float64(driftedAll) / float64(totalAll)
	}
	return RatioReport{Entries: entries, OverallRatio: overall}
}

// DriftRatioHasDrift returns true when any resource type has a non-zero drift ratio.
func DriftRatioHasDrift(rr RatioReport) bool {
	return rr.OverallRatio > 0
}

// FprintDriftRatio writes a human-readable drift ratio table to w.
func FprintDriftRatio(w io.Writer, rr RatioReport) {
	if len(rr.Entries) == 0 {
		fmt.Fprintln(w, "drift-ratio: no resources found")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %10s\n", "TYPE", "TOTAL", "DRIFTED", "RATIO")
	fmt.Fprintf(w, "%-30s %8s %8s %10s\n", "----", "-----", "-------", "-----")
	for _, e := range rr.Entries {
		fmt.Fprintf(w, "%-30s %8d %8d %9.1f%%\n", e.Type, e.Total, e.Drifted, e.RatioPct)
	}
	fmt.Fprintf(w, "\nOverall drift ratio: %.1f%%\n", rr.OverallRatio*100)
}
