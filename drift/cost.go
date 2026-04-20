package drift

import (
	"fmt"
	"io"
	"sort"
)

// CostWeight defines an estimated monthly cost weight per resource type.
var CostWeight = map[string]float64{
	"aws_instance":        50.0,
	"aws_rds_instance":    100.0,
	"aws_lambda_function": 5.0,
	"aws_s3_bucket":       2.0,
	"aws_elb":             20.0,
	"aws_nat_gateway":     35.0,
}

// CostEntry holds cost information for a single resource.
type CostEntry struct {
	ResourceType string  `json:"resource_type"`
	ResourceID   string  `json:"resource_id"`
	EstimatedUSD float64 `json:"estimated_usd"`
	Drifted      bool    `json:"drifted"`
}

// CostReport summarises estimated costs across a report.
type CostReport struct {
	Entries        []CostEntry `json:"entries"`
	TotalManaged   float64     `json:"total_managed_usd"`
	TotalUntracked float64     `json:"total_untracked_usd"`
}

// EstimateCosts builds a CostReport from a drift Report.
func EstimateCosts(r Report) CostReport {
	var cr CostReport
	for _, res := range r.Managed {
		w := costFor(res.Type)
		cr.Entries = append(cr.Entries, CostEntry{ResourceType: res.Type, ResourceID: res.ID, EstimatedUSD: w, Drifted: false})
		cr.TotalManaged += w
	}
	for _, res := range r.Untracked {
		w := costFor(res.Type)
		cr.Entries = append(cr.Entries, CostEntry{ResourceType: res.Type, ResourceID: res.ID, EstimatedUSD: w, Drifted: true})
		cr.TotalUntracked += w
	}
	sort.Slice(cr.Entries, func(i, j int) bool {
		return cr.Entries[i].EstimatedUSD > cr.Entries[j].EstimatedUSD
	})
	return cr
}

func costFor(resType string) float64 {
	if w, ok := CostWeight[resType]; ok {
		return w
	}
	return 1.0
}

// FprintCost writes a human-readable cost report to w.
func FprintCost(w io.Writer, cr CostReport) {
	fmt.Fprintln(w, "=== Cost Estimate ===")
	if len(cr.Entries) == 0 {
		fmt.Fprintln(w, "  No resources to estimate.")
		return
	}
	for _, e := range cr.Entries {
		tag := "managed"
		if e.Drifted {
			tag = "untracked"
		}
		fmt.Fprintf(w, "  %-40s %-12s $%7.2f/mo\n", e.ResourceID, "("+tag+")", e.EstimatedUSD)
	}
	fmt.Fprintf(w, "\n  Managed total:   $%.2f/mo\n", cr.TotalManaged)
	fmt.Fprintf(w, "  Untracked total: $%.2f/mo\n", cr.TotalUntracked)
	fmt.Fprintf(w, "  Grand total:     $%.2f/mo\n", cr.TotalManaged+cr.TotalUntracked)
}
