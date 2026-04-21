package drift

import (
	"fmt"
	"io"
	"sort"
)

// ImpactLevel represents the estimated blast radius of a drifted resource.
type ImpactLevel string

const (
	ImpactCritical ImpactLevel = "critical"
	ImpactHigh     ImpactLevel = "high"
	ImpactMedium   ImpactLevel = "medium"
	ImpactLow      ImpactLevel = "low"
)

// ImpactResult holds the evaluated impact for a single resource.
type ImpactResult struct {
	ResourceKey string
	Type        string
	Level       ImpactLevel
	Reason      string
}

// ImpactReport aggregates all impact results.
type ImpactReport struct {
	Results []ImpactResult
}

// impactRules maps resource types to their impact level.
var impactRules = map[string]ImpactLevel{
	"aws_iam_role":          ImpactCritical,
	"aws_iam_policy":        ImpactCritical,
	"aws_security_group":    ImpactHigh,
	"aws_vpc":               ImpactHigh,
	"aws_subnet":            ImpactMedium,
	"aws_s3_bucket":         ImpactHigh,
	"aws_lambda_function":   ImpactMedium,
	"aws_instance":          ImpactMedium,
	"aws_route_table":       ImpactMedium,
	"aws_cloudwatch_alarm":  ImpactLow,
}

// EvaluateImpact assigns an impact level to each drifted resource in the report.
func EvaluateImpact(report Report) ImpactReport {
	var results []ImpactResult

	drifted := append(report.Missing, report.Untracked...)
	drifted = append(drifted, report.Changed...)

	for _, r := range drifted {
		level, ok := impactRules[r.Type]
		if !ok {
			level = ImpactLow
		}
		results = append(results, ImpactResult{
			ResourceKey: resourceKey(r),
			Type:        r.Type,
			Level:       level,
			Reason:      impactReason(level),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return impactOrder(results[i].Level) > impactOrder(results[j].Level)
	})

	return ImpactReport{Results: results}
}

func impactReason(level ImpactLevel) string {
	switch level {
		case ImpactCritical:
			return "security or identity resource; drift may expose access control issues"
		case ImpactHigh:
			return "networking or storage resource; drift may affect availability"
		case ImpactMedium:
			return "compute or routing resource; drift may affect workloads"
		default:
			return "low-risk resource; drift is unlikely to cause immediate issues"
	}
}

func impactOrder(level ImpactLevel) int {
	switch level {
		case ImpactCritical: return 4
		case ImpactHigh:     return 3
		case ImpactMedium:   return 2
		default:             return 1
	}
}

// ImpactHasCritical returns true if any result is critical.
func ImpactHasCritical(r ImpactReport) bool {
	for _, res := range r.Results {
		if res.Level == ImpactCritical {
			return true
		}
	}
	return false
}

// FprintImpact writes a human-readable impact report to w.
func FprintImpact(w io.Writer, r ImpactReport) {
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "No drifted resources to evaluate.")
		return
	}
	fmt.Fprintln(w, "Impact Assessment:")
	for _, res := range r.Results {
		fmt.Fprintf(w, "  [%-8s] %s (%s)\n", res.Level, res.ResourceKey, res.Reason)
	}
}
