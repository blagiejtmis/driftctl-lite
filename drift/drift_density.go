package drift

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// DriftDensityEntry holds drift density metrics for a single resource type.
type DriftDensityEntry struct {
	Type        string  `json:"type"`
	Total       int     `json:"total"`
	Drifted     int     `json:"drifted"`
	Density     float64 `json:"density"` // drifted / total
	DensityPct  float64 `json:"density_pct"`
}

// DriftDensityReport is the full output of BuildDriftDensity.
type DriftDensityReport struct {
	Entries        []DriftDensityEntry `json:"entries"`
	OverallDensity float64             `json:"overall_density"`
}

// BuildDriftDensity computes per-type drift density from a Report.
func BuildDriftDensity(report Report) DriftDensityReport {
	type counts struct{ total, drifted int }
	buckets := map[string]*counts{}

	for _, r := range report.Managed {
		if _, ok := buckets[r.Type]; !ok {
			buckets[r.Type] = &counts{}
		}
		buckets[r.Type].total++
	}
	for _, r := range report.Missing {
		if _, ok := buckets[r.Type]; !ok {
			buckets[r.Type] = &counts{}
		}
		buckets[r.Type].total++
		buckets[r.Type].drifted++
	}
	for _, r := range report.Untracked {
		if _, ok := buckets[r.Type]; !ok {
			buckets[r.Type] = &counts{}
		}
		buckets[r.Type].total++
		buckets[r.Type].drifted++
	}

	var entries []DriftDensityEntry
	totalAll, driftedAll := 0, 0
	for typ, c := range buckets {
		density := 0.0
		if c.total > 0 {
			density = float64(c.drifted) / float64(c.total)
		}
		entries = append(entries, DriftDensityEntry{
			Type:       typ,
			Total:      c.total,
			Drifted:    c.drifted,
			Density:    density,
			DensityPct: density * 100,
		})
		totalAll += c.total
		driftedAll += c.drifted
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].DensityPct != entries[j].DensityPct {
			return entries[i].DensityPct > entries[j].DensityPct
		}
		return entries[i].Type < entries[j].Type
	})

	overall := 0.0
	if totalAll > 0 {
		overall = float64(driftedAll) / float64(totalAll) * 100
	}
	return DriftDensityReport{Entries: entries, OverallDensity: overall}
}

// DriftDensityHasDrift returns true if any entry has drifted resources.
func DriftDensityHasDrift(r DriftDensityReport) bool {
	return r.OverallDensity > 0
}

// FprintDensity writes a human-readable density table to w.
func FprintDensity(w io.Writer, r DriftDensityReport) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "No resources found for density analysis.")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %10s\n", "TYPE", "TOTAL", "DRIFTED", "DENSITY%")
	fmt.Fprintf(w, "%-30s %8s %8s %10s\n", "----", "-----", "-------", "--------")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "%-30s %8d %8d %9.1f%%\n", e.Type, e.Total, e.Drifted, e.DensityPct)
	}
	fmt.Fprintf(w, "\nOverall drift density: %.1f%%\n", r.OverallDensity)
}

// FprintDensityToStdout is a convenience wrapper around FprintDensity.
func FprintDensityToStdout(r DriftDensityReport) {
	FprintDensity(os.Stdout, r)
}
