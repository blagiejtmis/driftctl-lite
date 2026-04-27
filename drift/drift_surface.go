package drift

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// SurfaceEntry represents the attack/drift surface for a single resource type.
type SurfaceEntry struct {
	Type        string  `json:"type"`
	Total       int     `json:"total"`
	Drifted     int     `json:"drifted"`
	SurfacePct  float64 `json:"surface_pct"`
	RiskLevel   string  `json:"risk_level"`
}

// SurfaceReport holds all surface entries.
type SurfaceReport struct {
	Entries []SurfaceEntry `json:"entries"`
	Overall float64        `json:"overall_pct"`
}

// BuildDriftSurface computes the drift surface across resource types.
func BuildDriftSurface(report Report) SurfaceReport {
	type bucket struct{ total, drifted int }
	byType := map[string]*bucket{}

	for _, r := range report.Managed {
		if _, ok := byType[r.Type]; !ok {
			byType[r.Type] = &bucket{}
		}
		byType[r.Type].total++
	}
	for _, r := range report.Missing {
		if _, ok := byType[r.Type]; !ok {
			byType[r.Type] = &bucket{}
		}
		byType[r.Type].total++
		byType[r.Type].drifted++
	}
	for _, r := range report.Untracked {
		if _, ok := byType[r.Type]; !ok {
			byType[r.Type] = &bucket{}
		}
		byType[r.Type].total++
		byType[r.Type].drifted++
	}

	var entries []SurfaceEntry
	totalAll, driftedAll := 0, 0
	for t, b := range byType {
		pct := 0.0
		if b.total > 0 {
			pct = float64(b.drifted) / float64(b.total) * 100.0
		}
		entries = append(entries, SurfaceEntry{
			Type:       t,
			Total:      b.total,
			Drifted:    b.drifted,
			SurfacePct: pct,
			RiskLevel:  surfaceRisk(pct),
		})
		totalAll += b.total
		driftedAll += b.drifted
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].SurfacePct != entries[j].SurfacePct {
			return entries[i].SurfacePct > entries[j].SurfacePct
		}
		return entries[i].Type < entries[j].Type
	})

	overall := 0.0
	if totalAll > 0 {
		overall = float64(driftedAll) / float64(totalAll) * 100.0
	}
	return SurfaceReport{Entries: entries, Overall: overall}
}

func surfaceRisk(pct float64) string {
	switch {
	case pct >= 75:
		return "critical"
	case pct >= 50:
		return "high"
	case pct >= 25:
		return "medium"
	default:
		return "low"
	}
}

// SurfaceHasCritical returns true if any entry is critical risk.
func SurfaceHasCritical(r SurfaceReport) bool {
	for _, e := range r.Entries {
		if e.RiskLevel == "critical" {
			return true
		}
	}
	return false
}

// FprintSurface writes a human-readable surface report to w.
func FprintSurface(w io.Writer, r SurfaceReport) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "drift surface: no resources")
		return
	}
	fmt.Fprintf(w, "drift surface (overall: %.1f%%)\n", r.Overall)
	fmt.Fprintf(w, "  %-30s %6s %7s %8s %8s\n", "TYPE", "TOTAL", "DRIFTED", "PCT", "RISK")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-30s %6d %7d %7.1f%% %8s\n",
			e.Type, e.Total, e.Drifted, e.SurfacePct, e.RiskLevel)
	}
}

// FprintSurfaceToStdout is a convenience wrapper.
func FprintSurfaceToStdout(r SurfaceReport) { FprintSurface(os.Stdout, r) }
