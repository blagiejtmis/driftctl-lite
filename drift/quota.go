package drift

import (
	"fmt"
	"io"
	"sort"
)

// QuotaLimit defines a soft/hard limit for a resource type.
type QuotaLimit struct {
	Type      string
	SoftLimit int
	HardLimit int
}

// QuotaResult holds the evaluation result for a single resource type.
type QuotaResult struct {
	Type      string
	Count     int
	SoftLimit int
	HardLimit int
	Exceeded  string // "none", "soft", "hard"
}

// QuotaReport is the full output of EvaluateQuota.
type QuotaReport struct {
	Results []QuotaResult
}

// DefaultQuotaLimits returns a set of sensible built-in quota limits.
func DefaultQuotaLimits() []QuotaLimit {
	return []QuotaLimit{
		{Type: "aws_instance", SoftLimit: 20, HardLimit: 50},
		{Type: "aws_s3_bucket", SoftLimit: 50, HardLimit: 100},
		{Type: "aws_security_group", SoftLimit: 100, HardLimit: 200},
		{Type: "aws_iam_role", SoftLimit: 50, HardLimit: 100},
		{Type: "google_compute_instance", SoftLimit: 20, HardLimit: 50},
	}
}

// EvaluateQuota counts resources per type and checks against limits.
func EvaluateQuota(resources []Resource, limits []QuotaLimit) QuotaReport {
	counts := make(map[string]int)
	for _, r := range resources {
		counts[r.Type]++
	}

	limitMap := make(map[string]QuotaLimit)
	for _, l := range limits {
		limitMap[l.Type] = l
	}

	var results []QuotaResult
	for typ, count := range counts {
		l, ok := limitMap[typ]
		if !ok {
			continue
		}
		exceeded := "none"
		if count >= l.HardLimit {
			exceeded = "hard"
		} else if count >= l.SoftLimit {
			exceeded = "soft"
		}
		results = append(results, QuotaResult{
			Type:      typ,
			Count:     count,
			SoftLimit: l.SoftLimit,
			HardLimit: l.HardLimit,
			Exceeded:  exceeded,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Type < results[j].Type
	})

	return QuotaReport{Results: results}
}

// QuotaHasViolations returns true if any soft or hard limit is exceeded.
func QuotaHasViolations(r QuotaReport) bool {
	for _, res := range r.Results {
		if res.Exceeded != "none" {
			return true
		}
	}
	return false
}

// FprintQuota writes a human-readable quota report to w.
func FprintQuota(w io.Writer, r QuotaReport) {
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "[quota] no tracked resource types found")
		return
	}
	fmt.Fprintln(w, "[quota] resource quota evaluation:")
	for _, res := range r.Results {
		status := "OK"
		if res.Exceeded == "hard" {
			status = "HARD LIMIT EXCEEDED"
		} else if res.Exceeded == "soft" {
			status = "soft limit exceeded"
		}
		fmt.Fprintf(w, "  %-30s count=%-4d soft=%-4d hard=%-4d [%s]\n",
			res.Type, res.Count, res.SoftLimit, res.HardLimit, status)
	}
}
