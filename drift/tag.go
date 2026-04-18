package drift

import "fmt"

// TagRule defines a required tag key and optional allowed values.
type TagRule struct {
	Key    string   `json:"key"`
	Values []string `json:"values,omitempty"`
}

// TagPolicy holds a set of tag rules to enforce.
type TagPolicy struct {
	Required []TagRule `json:"required"`
}

// TagViolation describes a single tag compliance failure.
type TagViolation struct {
	Resource Resource
	Rule     TagRule
	Reason   string
}

// EvaluateTags checks resources against a TagPolicy and returns violations.
func EvaluateTags(resources []Resource, policy TagPolicy) []TagViolation {
	var violations []TagViolation
	for _, res := range resources {
		for _, rule := range policy.Required {
			val, ok := res.Attributes[rule.Key]
			if !ok {
				violations = append(violations, TagViolation{
					Resource: res,
					Rule:     rule,
					Reason:   fmt.Sprintf("missing required tag %q", rule.Key),
				})
				continue
			}
			if len(rule.Values) > 0 && !containsStr(rule.Values, fmt.Sprintf("%v", val)) {
				violations = append(violations, TagViolation{
					Resource: res,
					Rule:     rule,
					Reason:   fmt.Sprintf("tag %q value %q not in allowed list", rule.Key, val),
				})
			}
		}
	}
	return violations
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
