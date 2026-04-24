package drift

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// DriftMapEntry holds per-resource-type drift statistics.
type DriftMapEntry struct {
	Type      string
	Total     int
	Missing   int
	Untracked int
	Changed   int
	DriftPct  float64
}

// DriftMap is an ordered slice of DriftMapEntry values.
type DriftMap []DriftMapEntry

// BuildDriftMap aggregates drift statistics by resource type from a Report.
func BuildDriftMap(r Report) DriftMap {
	type counts struct {
		total, missing, untracked, changed int
	}
	m := map[string]*counts{}

	ensure := func(t string) {
		if _, ok := m[t]; !ok {
			m[t] = &counts{}
		}
	}

	for _, res := range r.Managed {
		ensure(res.Type)
		m[res.Type].total++
	}
	for _, res := range r.Missing {
		ensure(res.Type)
		m[res.Type].total++
		m[res.Type].missing++
	}
	for _, res := range r.Untracked {
		ensure(res.Type)
		m[res.Type].total++
		m[res.Type].untracked++
	}
	for _, res := range r.Changed {
		ensure(res.Type)
		m[res.Type].total++
		m[res.Type].changed++
	}

	result := make(DriftMap, 0, len(m))
	for t, c := range m {
		drifted := c.missing + c.untracked + c.changed
		pct := 0.0
		if c.total > 0 {
			pct = float64(drifted) / float64(c.total) * 100.0
		}
		result = append(result, DriftMapEntry{
			Type:      t,
			Total:     c.total,
			Missing:   c.missing,
			Untracked: c.untracked,
			Changed:   c.changed,
			DriftPct:  pct,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DriftPct != result[j].DriftPct {
			return result[i].DriftPct > result[j].DriftPct
		}
		return result[i].Type < result[j].Type
	})
	return result
}

// DriftMapHasDrift returns true if any entry has drift.
func DriftMapHasDrift(dm DriftMap) bool {
	for _, e := range dm {
		if e.Missing+e.Untracked+e.Changed > 0 {
			return true
		}
	}
	return false
}

// FprintDriftMap writes a formatted drift map table to w.
func FprintDriftMap(w io.Writer, dm DriftMap) {
	if len(dm) == 0 {
		fmt.Fprintln(w, "drift-map: no resources found")
		return
	}
	fmt.Fprintf(w, "%-30s %6s %7s %9s %7s %8s\n",
		"TYPE", "TOTAL", "MISSING", "UNTRACKED", "CHANGED", "DRIFT%")
	fmt.Fprintln(w, strings.Repeat("-", 72))
	for _, e := range dm {
		fmt.Fprintf(w, "%-30s %6d %7d %9d %7d %7.1f%%\n",
			e.Type, e.Total, e.Missing, e.Untracked, e.Changed, e.DriftPct)
	}
}
