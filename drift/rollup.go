package drift

import (
	"fmt"
	"io"
	"sort"
)

// RollupEntry summarises drift counts for a single resource type.
type RollupEntry struct {
	Type      string `json:"type"`
	Managed   int    `json:"managed"`
	Missing   int    `json:"missing"`
	Untracked int    `json:"untracked"`
	Changed   int    `json:"changed"`
	Total     int    `json:"total"`
}

// Rollup aggregates a Report into per-type drift summaries.
func Rollup(r Report) []RollupEntry {
	counts := map[string]*RollupEntry{}

	ensure := func(t string) *RollupEntry {
		if _, ok := counts[t]; !ok {
			counts[t] = &RollupEntry{Type: t}
		}
		return counts[t]
	}

	for _, res := range r.Managed {
		ensure(res.Type).Managed++
	}
	for _, res := range r.Missing {
		ensure(res.Type).Missing++
	}
	for _, res := range r.Untracked {
		ensure(res.Type).Untracked++
	}
	for _, res := range r.Changed {
		ensure(res.Type).Changed++
	}

	entries := make([]RollupEntry, 0, len(counts))
	for _, e := range counts {
		e.Total = e.Managed + e.Missing + e.Untracked + e.Changed
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Type < entries[j].Type
	})
	return entries
}

// FprintRollup writes a human-readable rollup table to w.
func FprintRollup(w io.Writer, entries []RollupEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No resources to roll up.")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %10s %8s %8s\n",
		"TYPE", "MANAGED", "MISSING", "UNTRACKED", "CHANGED", "TOTAL")
	fmt.Fprintf(w, "%s\n", "----------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Fprintf(w, "%-30s %8d %8d %10d %8d %8d\n",
			e.Type, e.Managed, e.Missing, e.Untracked, e.Changed, e.Total)
	}
}

// RollupHasDrift returns true if any entry has drift.
func RollupHasDrift(entries []RollupEntry) bool {
	for _, e := range entries {
		if e.Missing > 0 || e.Untracked > 0 || e.Changed > 0 {
			return true
		}
	}
	return false
}
