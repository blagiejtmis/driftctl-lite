package drift

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// HotspotEntry represents a resource type with repeated drift occurrences.
type HotspotEntry struct {
	Type       string
	DriftCount int
	Total      int
	Score      float64 // DriftCount / Total * 100
}

// HotspotReport holds all hotspot entries.
type HotspotReport struct {
	Entries []HotspotEntry
}

// BuildHotspot identifies resource types with the highest drift frequency.
func BuildHotspot(reports []Report) HotspotReport {
	type counts struct {
		drifted int
		total   int
	}
	agg := map[string]*counts{}

	for _, r := range reports {
		for _, res := range r.Managed {
			if _, ok := agg[res.Type]; !ok {
				agg[res.Type] = &counts{}
			}
			agg[res.Type].total++
		}
		for _, res := range r.Missing {
			if _, ok := agg[res.Type]; !ok {
				agg[res.Type] = &counts{}
			}
			agg[res.Type].total++
			agg[res.Type].drifted++
		}
		for _, res := range r.Untracked {
			if _, ok := agg[res.Type]; !ok {
				agg[res.Type] = &counts{}
			}
			agg[res.Type].total++
			agg[res.Type].drifted++
		}
	}

	var entries []HotspotEntry
	for typ, c := range agg {
		score := 0.0
		if c.total > 0 {
			score = float64(c.drifted) / float64(c.total) * 100
		}
		entries = append(entries, HotspotEntry{
			Type:       typ,
			DriftCount: c.drifted,
			Total:      c.total,
			Score:      score,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Type < entries[j].Type
	})

	return HotspotReport{Entries: entries}
}

// HotspotHasEntries returns true when at least one drifted entry exists.
func HotspotHasEntries(r HotspotReport) bool {
	for _, e := range r.Entries {
		if e.DriftCount > 0 {
			return true
		}
	}
	return false
}

// FprintHotspot writes a formatted hotspot table to w.
func FprintHotspot(w io.Writer, r HotspotReport) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "No hotspot data available.")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %8s\n", "TYPE", "DRIFTED", "TOTAL", "SCORE")
	fmt.Fprintf(w, "%-30s %8s %8s %8s\n", "----", "-------", "-----", "-----")
	for _, e := range r.Entries {
		if e.DriftCount == 0 {
			continue
		}
		fmt.Fprintf(w, "%-30s %8d %8d %7.1f%%\n", e.Type, e.DriftCount, e.Total, e.Score)
	}
}

// FprintHotspotToStdout writes hotspot output to os.Stdout.
func FprintHotspotToStdout(r HotspotReport) {
	FprintHotspot(os.Stdout, r)
}
