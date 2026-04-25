package drift

import (
	"fmt"
	"io"
	"sort"
)

// BudgetConfig defines thresholds for allowed drift counts per resource type.
type BudgetConfig struct {
	MaxMissing   int            `json:"max_missing"`
	MaxUntracked int            `json:"max_untracked"`
	PerType      map[string]int `json:"per_type,omitempty"`
}

// BudgetResult holds the evaluation result for a single resource type.
type BudgetResult struct {
	Type      string
	DriftCount int
	Budget    int
	Exceeded  bool
}

// BudgetReport is the full output of EvaluateBudget.
type BudgetReport struct {
	Config  BudgetConfig
	Results []BudgetResult
	TotalMissing   int
	TotalUntracked int
	MissingExceeded   bool
	UntrackedExceeded bool
}

// DefaultBudgetConfig returns a sensible default budget.
func DefaultBudgetConfig() BudgetConfig {
	return BudgetConfig{
		MaxMissing:   5,
		MaxUntracked: 10,
		PerType:      map[string]int{},
	}
}

// EvaluateBudget checks drift counts against configured budgets.
func EvaluateBudget(report Report, cfg BudgetConfig) BudgetReport {
	typeCounts := map[string]int{}
	for _, r := range report.Missing {
		typeCounts[r.Type]++
	}
	for _, r := range report.Untracked {
		typeCounts[r.Type]++
	}

	var results []BudgetResult
	for t, count := range typeCounts {
		budget, ok := cfg.PerType[t]
		if !ok {
			continue
		}
		results = append(results, BudgetResult{
			Type:       t,
			DriftCount: count,
			Budget:     budget,
			Exceeded:   count > budget,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Type < results[j].Type
	})

	totalMissing := len(report.Missing)
	totalUntracked := len(report.Untracked)
	return BudgetReport{
		Config:            cfg,
		Results:           results,
		TotalMissing:      totalMissing,
		TotalUntracked:    totalUntracked,
		MissingExceeded:   totalMissing > cfg.MaxMissing,
		UntrackedExceeded: totalUntracked > cfg.MaxUntracked,
	}
}

// BudgetHasViolations returns true if any budget threshold is exceeded.
func BudgetHasViolations(br BudgetReport) bool {
	if br.MissingExceeded || br.UntrackedExceeded {
		return true
	}
	for _, r := range br.Results {
		if r.Exceeded {
			return true
		}
	}
	return false
}

// FprintBudget writes a human-readable budget report to w.
func FprintBudget(w io.Writer, br BudgetReport) {
	fmt.Fprintf(w, "Drift Budget\n")
	fmt.Fprintf(w, "  Missing:   %d / %d", br.TotalMissing, br.Config.MaxMissing)
	if br.MissingExceeded {
		fmt.Fprintf(w, " [EXCEEDED]")
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Untracked: %d / %d", br.TotalUntracked, br.Config.MaxUntracked)
	if br.UntrackedExceeded {
		fmt.Fprintf(w, " [EXCEEDED]")
	}
	fmt.Fprintln(w)
	if len(br.Results) > 0 {
		fmt.Fprintln(w, "  Per-type budgets:")
		for _, r := range br.Results {
			status := "ok"
			if r.Exceeded {
				status = "EXCEEDED"
			}
			fmt.Fprintf(w, "    %-30s %d / %d [%s]\n", r.Type, r.DriftCount, r.Budget, status)
		}
	}
}
