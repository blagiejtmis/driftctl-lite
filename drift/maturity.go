package drift

import (
	"fmt"
	"io"
	"sort"
)

// MaturityLevel represents the IaC maturity tier.
type MaturityLevel string

const (
	MaturityInitial    MaturityLevel = "Initial"
	MaturityDeveloping MaturityLevel = "Developing"
	MaturityDefined    MaturityLevel = "Defined"
	MaturityManaged    MaturityLevel = "Managed"
	MaturityOptimizing MaturityLevel = "Optimizing"
)

// MaturityResult holds per-type maturity data.
type MaturityResult struct {
	Type       string
	Total      int
	Managed    int
	Coverage   float64
	Level      MaturityLevel
}

// MaturityReport is the full report across all resource types.
type MaturityReport struct {
	Results  []MaturityResult
	Overall  MaturityLevel
	AvgScore float64
}

// EvaluateMaturity computes maturity levels per resource type from a drift Report.
func EvaluateMaturity(report Report) MaturityReport {
	type bucket struct {
		total   int
		managed int
	}
	buckets := map[string]*bucket{}

	for _, r := range report.Managed {
		if _, ok := buckets[r.Type]; !ok {
			buckets[r.Type] = &bucket{}
		}
		buckets[r.Type].total++
		buckets[r.Type].managed++
	}
	for _, r := range report.Missing {
		if _, ok := buckets[r.Type]; !ok {
			buckets[r.Type] = &bucket{}
		}
		buckets[r.Type].total++
	}
	for _, r := range report.Untracked {
		if _, ok := buckets[r.Type]; !ok {
			buckets[r.Type] = &bucket{}
		}
		buckets[r.Type].total++
	}

	var results []MaturityResult
	var scoreSum float64

	for typ, b := range buckets {
		cov := 0.0
		if b.total > 0 {
			cov = float64(b.managed) / float64(b.total) * 100.0
		}
		lvl := maturityLevel(cov)
		results = append(results, MaturityResult{
			Type:     typ,
			Total:    b.total,
			Managed:  b.managed,
			Coverage: cov,
			Level:    lvl,
		})
		scoreSum += cov
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Coverage != results[j].Coverage {
			return results[i].Coverage > results[j].Coverage
		}
		return results[i].Type < results[j].Type
	})

	avg := 0.0
	if len(results) > 0 {
		avg = scoreSum / float64(len(results))
	}

	return MaturityReport{
		Results:  results,
		Overall:  maturityLevel(avg),
		AvgScore: avg,
	}
}

func maturityLevel(coverage float64) MaturityLevel {
	switch {
	case coverage >= 90:
		return MaturityOptimizing
	case coverage >= 75:
		return MaturityManaged
	case coverage >= 50:
		return MaturityDefined
	case coverage >= 25:
		return MaturityDeveloping
	default:
		return MaturityInitial
	}
}

// MaturityHasCritical returns true when overall maturity is Initial or Developing.
func MaturityHasCritical(r MaturityReport) bool {
	return r.Overall == MaturityInitial || r.Overall == MaturityDeveloping
}

// FprintMaturity writes a human-readable maturity report to w.
func FprintMaturity(w io.Writer, r MaturityReport) {
	fmt.Fprintf(w, "IaC Maturity Report\n")
	fmt.Fprintf(w, "Overall: %s (avg coverage %.1f%%)\n\n", r.Overall, r.AvgScore)
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "  No resource types found.")
		return
	}
	fmt.Fprintf(w, "  %-30s %8s %8s %8s %s\n", "Type", "Total", "Managed", "Cov%", "Level")
	for _, res := range r.Results {
		fmt.Fprintf(w, "  %-30s %8d %8d %7.1f%% %s\n",
			res.Type, res.Total, res.Managed, res.Coverage, res.Level)
	}
}
