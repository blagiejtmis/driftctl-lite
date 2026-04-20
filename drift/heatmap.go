package drift

import (
	"fmt"
	"io"
	"sort"
)

// HeatmapEntry represents drift intensity for a resource type.
type HeatmapEntry struct {
	Type     string `json:"type"`
	Total    int    `json:"total"`
	Drifted  int    `json:"drifted"`
	Intensity float64 `json:"intensity"` // 0.0 – 1.0
}

// HeatmapResult holds all entries for a heatmap report.
type HeatmapResult struct {
	Entries []HeatmapEntry `json:"entries"`
}

// BuildHeatmap computes per-type drift intensity from a Report.
func BuildHeatmap(r Report) HeatmapResult {
	type counts struct{ total, drifted int }
	m := map[string]*counts{}

	for _, res := range r.Managed {
		if m[res.Type] == nil {
			m[res.Type] = &counts{}
		}
		m[res.Type].total++
	}
	for _, res := range r.Missing {
		if m[res.Type] == nil {
			m[res.Type] = &counts{}
		}
		m[res.Type].total++
		m[res.Type].drifted++
	}
	for _, res := range r.Untracked {
		if m[res.Type] == nil {
			m[res.Type] = &counts{}
		}
		m[res.Type].total++
		m[res.Type].drifted++
	}

	entries := make([]HeatmapEntry, 0, len(m))
	for typ, c := range m {
		intensity := 0.0
		if c.total > 0 {
			intensity = float64(c.drifted) / float64(c.total)
		}
		entries = append(entries, HeatmapEntry{
			Type:      typ,
			Total:     c.total,
			Drifted:   c.drifted,
			Intensity: intensity,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Intensity != entries[j].Intensity {
			return entries[i].Intensity > entries[j].Intensity
		}
		return entries[i].Type < entries[j].Type
	})
	return HeatmapResult{Entries: entries}
}

// FprintHeatmap writes a human-readable heatmap to w.
func FprintHeatmap(w io.Writer, h HeatmapResult) {
	if len(h.Entries) == 0 {
		fmt.Fprintln(w, "No resources found for heatmap.")
		return
	}
	fmt.Fprintln(w, "Drift Heatmap (by resource type):")
	fmt.Fprintf(w, "  %-30s %6s %7s %8s\n", "TYPE", "TOTAL", "DRIFTED", "INTENSITY")
	for _, e := range h.Entries {
		bar := heatBar(e.Intensity)
		fmt.Fprintf(w, "  %-30s %6d %7d %8.0f%%  %s\n",
			e.Type, e.Total, e.Drifted, e.Intensity*100, bar)
	}
}

func heatBar(intensity float64) string {
	const width = 10
	filled := int(intensity * width)
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
