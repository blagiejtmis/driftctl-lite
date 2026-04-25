package drift

import (
	"fmt"
	"io"
	"sort"
)

// CoverageResult holds coverage statistics for a single resource type.
type CoverageResult struct {
	Type    string
	Total   int
	Managed int
	Pct     float64
}

// CoverageReport holds the full coverage analysis.
type CoverageReport struct {
	Results  []CoverageResult
	Overall  float64
	Total    int
	Managed  int
}

// EvaluateCoverage computes per-type and overall IaC coverage from a Report.
func EvaluateCoverage(report Report) CoverageReport {
	typeTotal := map[string]int{}
	typeManaged := map[string]int{}

	for _, r := range report.Managed {
		typeTotal[r.Type]++
		typeManaged[r.Type]++
	}
	for _, r := range report.Missing {
		typeTotal[r.Type]++
	}
	for _, r := range report.Untracked {
		typeTotal[r.Type]++
	}

	var results []CoverageResult
	for t, total := range typeTotal {
		managed := typeManaged[t]
		pct := 0.0
		if total > 0 {
			pct = float64(managed) / float64(total) * 100.0
		}
		results = append(results, CoverageResult{
			Type:    t,
			Total:   total,
			Managed: managed,
			Pct:     pct,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Pct != results[j].Pct {
			return results[i].Pct < results[j].Pct
		}
		return results[i].Type < results[j].Type
	})

	totalAll := len(report.Managed) + len(report.Missing) + len(report.Untracked)
	managedAll := len(report.Managed)
	overall := 0.0
	if totalAll > 0 {
		overall = float64(managedAll) / float64(totalAll) * 100.0
	}

	return CoverageReport{
		Results: results,
		Overall: overall,
		Total:   totalAll,
		Managed: managedAll,
	}
}

// CoverageHasGaps returns true if any resource type has less than 100% coverage.
func CoverageHasGaps(cr CoverageReport) bool {
	return cr.Overall < 100.0 && cr.Total > 0
}

// FprintCoverage writes a human-readable coverage report to w.
func FprintCoverage(w io.Writer, cr CoverageReport) {
	if cr.Total == 0 {
		fmt.Fprintln(w, "coverage: no resources found")
		return
	}
	fmt.Fprintf(w, "IaC Coverage Report\n")
	fmt.Fprintf(w, "Overall: %.1f%% (%d/%d managed)\n\n", cr.Overall, cr.Managed, cr.Total)
	for _, r := range cr.Results {
		bar := coverageBar(r.Pct)
		fmt.Fprintf(w, "  %-30s %s %.1f%% (%d/%d)\n", r.Type, bar, r.Pct, r.Managed, r.Total)
	}
}

func coverageBar(pct float64) string {
	filled := int(pct / 10)
	bar := ""
	for i := 0; i < 10; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
