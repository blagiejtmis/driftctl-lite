package drift

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// IndexEntry holds a summary row for a single resource type in the drift index.
type IndexEntry struct {
	Type        string
	Total       int
	Managed     int
	Missing     int
	Untracked   int
	Changed     int
	DriftPct    float64
}

// DriftIndex maps resource type to its IndexEntry.
type DriftIndex struct {
	Entries []IndexEntry
	Total   IndexEntry // aggregate row
}

// BuildDriftIndex constructs a DriftIndex from a Report.
func BuildDriftIndex(report Report) DriftIndex {
	byType := map[string]*IndexEntry{}

	for _, r := range report.Managed {
		e := getOrCreate(byType, r.Type)
		e.Total++
		e.Managed++
	}
	for _, r := range report.Missing {
		e := getOrCreate(byType, r.Type)
		e.Total++
		e.Missing++
	}
	for _, r := range report.Untracked {
		e := getOrCreate(byType, r.Type)
		e.Total++
		e.Untracked++
	}
	for _, r := range report.Changed {
		e := getOrCreate(byType, r.Type)
		e.Changed++
	}

	entries := make([]IndexEntry, 0, len(byType))
	agg := IndexEntry{Type: "(all)"}
	for _, e := range byType {
		drifted := e.Missing + e.Untracked + e.Changed
		if e.Total > 0 {
			e.DriftPct = float64(drifted) / float64(e.Total) * 100
		}
		agg.Total += e.Total
		agg.Managed += e.Managed
		agg.Missing += e.Missing
		agg.Untracked += e.Untracked
		agg.Changed += e.Changed
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].DriftPct != entries[j].DriftPct {
			return entries[i].DriftPct > entries[j].DriftPct
		}
		return entries[i].Type < entries[j].Type
	})
	driftedAgg := agg.Missing + agg.Untracked + agg.Changed
	if agg.Total > 0 {
		agg.DriftPct = float64(driftedAgg) / float64(agg.Total) * 100
	}
	return DriftIndex{Entries: entries, Total: agg}
}

// DriftIndexHasDrift returns true if any entry has drift.
func DriftIndexHasDrift(idx DriftIndex) bool {
	return idx.Total.Missing+idx.Total.Untracked+idx.Total.Changed > 0
}

// FprintDriftIndex writes a human-readable drift index table to w.
func FprintDriftIndex(w io.Writer, idx DriftIndex) {
	if len(idx.Entries) == 0 {
		fmt.Fprintln(w, "drift-index: no resources found")
		return
	}
	fmt.Fprintf(w, "%-30s %6s %7s %7s %9s %7s %8s\n",
		"TYPE", "TOTAL", "MANAGED", "MISSING", "UNTRACKED", "CHANGED", "DRIFT%")
	fmt.Fprintln(w, strings.Repeat("-", 80))
	for _, e := range idx.Entries {
		fmt.Fprintf(w, "%-30s %6d %7d %7d %9d %7d %7.1f%%\n",
			e.Type, e.Total, e.Managed, e.Missing, e.Untracked, e.Changed, e.DriftPct)
	}
	fmt.Fprintln(w, strings.Repeat("-", 80))
	fmt.Fprintf(w, "%-30s %6d %7d %7d %9d %7d %7.1f%%\n",
		idx.Total.Type, idx.Total.Total, idx.Total.Managed, idx.Total.Missing,
		idx.Total.Untracked, idx.Total.Changed, idx.Total.DriftPct)
}

func getOrCreate(m map[string]*IndexEntry, t string) *IndexEntry {
	if _, ok := m[t]; !ok {
		m[t] = &IndexEntry{Type: t}
	}
	return m[t]
}
