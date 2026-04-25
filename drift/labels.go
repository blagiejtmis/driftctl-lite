package drift

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// LabelRule defines a rule for required or disallowed labels.
type LabelRule struct {
	Key      string   `json:"key"`
	Required bool     `json:"required"`
	Allowed  []string `json:"allowed,omitempty"`
}

// LabelViolation records a label policy violation for a resource.
type LabelViolation struct {
	Resource Resource
	Rule     LabelRule
	Message  string
}

// EvaluateLabels checks resources against label rules and returns violations.
func EvaluateLabels(resources []Resource, rules []LabelRule) []LabelViolation {
	var violations []LabelViolation
	for _, res := range resources {
		for _, rule := range rules {
			val, exists := res.Attributes[rule.Key]
			if rule.Required && !exists {
				violations = append(violations, LabelViolation{
					Resource: res,
					Rule:     rule,
					Message:  fmt.Sprintf("missing required label %q", rule.Key),
				})
				continue
			}
			if exists && len(rule.Allowed) > 0 && !containsStr(rule.Allowed, fmt.Sprintf("%v", val)) {
				violations = append(violations, LabelViolation{
					Resource: res,
					Rule:     rule,
					Message:  fmt.Sprintf("label %q has disallowed value %q", rule.Key, val),
				})
			}
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		ki := resourceKey(violations[i].Resource)
		kj := resourceKey(violations[j].Resource)
		if ki != kj {
			return ki < kj
		}
		return violations[i].Rule.Key < violations[j].Rule.Key
	})
	return violations
}

// LabelHasViolations returns true when there is at least one violation.
func LabelHasViolations(violations []LabelViolation) bool {
	return len(violations) > 0
}

// FprintLabels writes a human-readable label violation report.
func FprintLabels(w io.Writer, violations []LabelViolation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "[labels] no violations found")
		return
	}
	fmt.Fprintf(w, "[labels] %d violation(s) found:\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(w, "  %-12s %-30s %s\n",
			v.Resource.Type,
			v.Resource.ID,
			strings.TrimSpace(v.Message),
		)
	}
}
