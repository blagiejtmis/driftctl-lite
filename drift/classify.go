package drift

import "fmt"

// Classification severity levels
const (
	ClassLow      = "low"
	ClassMedium   = "medium"
	ClassHigh     = "high"
	ClassCritical = "critical"
)

// ClassifyResult holds the classification for a single resource.
type ClassifyResult struct {
	Resource Resource
	Severity string
	Reason   string
}

// ClassifyReport holds all classification results.
type ClassifyReport struct {
	Results []ClassifyResult
}

// severityOrder maps severity to sort weight.
var severityOrder = map[string]int{
	ClassCritical: 4,
	ClassHigh:     3,
	ClassMedium:   2,
	ClassLow:      1,
}

// classifyType returns a severity level based on resource type.
func classifyType(rtype string) (string, string) {
	switch rtype {
	case "aws_iam_role", "aws_iam_policy", "aws_iam_user":
		return ClassCritical, "IAM resources carry privilege escalation risk"
	case "aws_security_group", "aws_network_acl":
		return ClassHigh, "Network controls affect blast radius"
	case "aws_s3_bucket", "aws_rds_instance":
		return ClassMedium, "Data-plane resource with potential data exposure"
	default:
		return ClassLow, "General resource with limited security impact"
	}
}

// ClassifyHasCritical returns true if any result is critical or high.
func ClassifyHasCritical(r ClassifyReport) bool {
	for _, res := range r.Results {
		if res.Severity == ClassCritical || res.Severity == ClassHigh {
			return true
		}
	}
	return false
}

// Classify assigns severity classifications to drifted resources.
func Classify(report Report) ClassifyReport {
	var results []ClassifyResult
	for _, res := range report.Missing {
		sev, reason := classifyType(res.Type)
		results = append(results, ClassifyResult{Resource: res, Severity: sev, Reason: fmt.Sprintf("missing: %s", reason)})
	}
	for _, res := range report.Untracked {
		sev, reason := classifyType(res.Type)
		results = append(results, ClassifyResult{Resource: res, Severity: sev, Reason: fmt.Sprintf("untracked: %s", reason)})
	}
	// sort descending by severity
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && severityOrder[results[j].Severity] > severityOrder[results[j-1].Severity]; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return ClassifyReport{Results: results}
}
