package drift

import (
	"fmt"
	"io"
	"sort"
)

// ComplianceFramework represents a named set of required resource types.
type ComplianceFramework struct {
	Name     string   `json:"name"`
	Required []string `json:"required_types"`
}

// ComplianceResult holds the evaluation outcome for a single framework.
type ComplianceResult struct {
	Framework string
	Missing   []string
	Passing   bool
}

// EvaluateCompliance checks whether all required resource types from each
// framework are present in the managed resources of the report.
func EvaluateCompliance(report Report, frameworks []ComplianceFramework) []ComplianceResult {
	present := map[string]bool{}
	for _, r := range report.Managed {
		present[r.Type] = true
	}

	results := make([]ComplianceResult, 0, len(frameworks))
	for _, fw := range frameworks {
		var missing []string
		for _, req := range fw.Required {
			if !present[req] {
				missing = append(missing, req)
			}
		}
		sort.Strings(missing)
		results = append(results, ComplianceResult{
			Framework: fw.Name,
			Missing:   missing,
			Passing:   len(missing) == 0,
		})
	}
	return results
}

// ComplianceHasFailures returns true if any framework result is not passing.
func ComplianceHasFailures(results []ComplianceResult) bool {
	for _, r := range results {
		if !r.Passing {
			return true
		}
	}
	return false
}

// FprintCompliance writes a human-readable compliance report to w.
func FprintCompliance(w io.Writer, results []ComplianceResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No compliance frameworks evaluated.")
		return
	}
	for _, r := range results {
		status := "PASS"
		if !r.Passing {
			status = "FAIL"
		}
		fmt.Fprintf(w, "[%s] %s\n", status, r.Framework)
		for _, m := range r.Missing {
			fmt.Fprintf(w, "  - missing type: %s\n", m)
		}
	}
}
