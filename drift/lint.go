package drift

import "fmt"

// LintRule represents a single lint rule with a name and description.
type LintRule struct {
	Name    string
	Message string
}

// LintResult holds a resource and the rules it violated.
type LintResult struct {
	Resource Resource
	Violations []LintRule
}

// LintReport holds all lint results.
type LintReport struct {
	Results []LintResult
}

// Lint checks resources against basic lint rules and returns a report.
func Lint(resources []Resource) LintReport {
	var report LintReport
	for _, r := range resources {
		var violations []LintRule
		if r.ID == "" {
			violations = append(violations, LintRule{
				Name:    "missing-id",
				Message: "resource has an empty ID",
			})
		}
		if r.Type == "" {
			violations = append(violations, LintRule{
				Name:    "missing-type",
				Message: "resource has an empty type",
			})
		}
		if len(r.Attributes) == 0 {
			violations = append(violations, LintRule{
				Name:    "no-attributes",
				Message: "resource has no attributes defined",
			})
		}
		if len(violations) > 0 {
			report.Results = append(report.Results, LintResult{
				Resource:   r,
				Violations: violations,
			})
		}
	}
	return report
}

// LintHasErrors returns true if any lint violations exist.
func LintHasErrors(r LintReport) bool {
	return len(r.Results) > 0
}

// FprintLint writes a human-readable lint report to w.
func FprintLint(w interface{ WriteString(string) (int, error) }, r LintReport) {
	if !LintHasErrors(r) {
		w.WriteString("lint: no issues found\n")
		return
	}
	w.WriteString(fmt.Sprintf("lint: %d resource(s) with issues\n", len(r.Results)))
	for _, res := range r.Results {
		w.WriteString(fmt.Sprintf("  [%s] %s\n", res.Resource.Type, res.Resource.ID))
		for _, v := range res.Violations {
			w.WriteString(fmt.Sprintf("    - %s: %s\n", v.Name, v.Message))
		}
	}
}
