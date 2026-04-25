package drift

import (
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

// VelocityEntry represents the drift rate for a resource type over a window.
type VelocityEntry struct {
	Type          string  `json:"type"`
	DriftedCount  int     `json:"drifted_count"`
	TotalCount    int     `json:"total_count"`
	DriftRate     float64 `json:"drift_rate"`
	DeltaFromPrev float64 `json:"delta_from_prev"`
	Trend         string  `json:"trend"` // "increasing", "decreasing", "stable"
}

// VelocityReport holds the computed drift velocity results.
type VelocityReport struct {
	ComputedAt time.Time       `json:"computed_at"`
	WindowDays int             `json:"window_days"`
	Entries    []VelocityEntry `json:"entries"`
}

// VelocityHasIncreasing returns true if any entry has an increasing drift trend.
func VelocityHasIncreasing(r VelocityReport) bool {
	for _, e := range r.Entries {
		if e.Trend == "increasing" {
			return true
		}
	}
	return false
}

// EvaluateVelocity computes drift velocity per resource type using two snapshots.
// current represents the latest drift state; previous is an earlier snapshot.
func EvaluateVelocity(current, previous Report, windowDays int) VelocityReport {
	if windowDays <= 0 {
		windowDays = 7
	}

	curByType := make(map[string][2]int) // [drifted, total]
	for _, r := range current.Missing {
		e := curByType[r.Type]
		e[0]++
		e[1]++
		curByType[r.Type] = e
	}
	for _, r := range current.Untracked {
		e := curByType[r.Type]
		e[0]++
		e[1]++
		curByType[r.Type] = e
	}
	for _, r := range current.Managed {
		e := curByType[r.Type]
		e[1]++
		curByType[r.Type] = e
	}

	prevByType := make(map[string]float64)
	for _, r := range previous.Missing {
		prevByType[r.Type]++
	}
	for _, r := range previous.Untracked {
		prevByType[r.Type]++
	}

	var entries []VelocityEntry
	for t, counts := range curByType {
		drifted := counts[0]
		total := counts[1]
		var rate float64
		if total > 0 {
			rate = math.Round((float64(drifted)/float64(total))*10000) / 100
		}
		prevDrifted := prevByType[t]
		delta := float64(drifted) - prevDrifted
		trend := "stable"
		if delta > 0 {
			trend = "increasing"
		} else if delta < 0 {
			trend = "decreasing"
		}
		entries = append(entries, VelocityEntry{
			Type:          t,
			DriftedCount:  drifted,
			TotalCount:    total,
			DriftRate:     rate,
			DeltaFromPrev: delta,
			Trend:         trend,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].DriftRate != entries[j].DriftRate {
			return entries[i].DriftRate > entries[j].DriftRate
		}
		return entries[i].Type < entries[j].Type
	})

	return VelocityReport{
		ComputedAt: time.Now().UTC(),
		WindowDays: windowDays,
		Entries:    entries,
	}
}

// FprintVelocity writes a human-readable velocity report to w.
func FprintVelocity(w io.Writer, r VelocityReport) {
	fmt.Fprintf(w, "Drift Velocity Report (window: %d days)\n", r.WindowDays)
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "  No resources tracked.")
		return
	}
	fmt.Fprintf(w, "  %-25s %8s %8s %10s %8s\n", "TYPE", "DRIFTED", "TOTAL", "RATE(%)", "TREND")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-25s %8d %8d %10.2f %8s\n",
			e.Type, e.DriftedCount, e.TotalCount, e.DriftRate, e.Trend)
	}
}
