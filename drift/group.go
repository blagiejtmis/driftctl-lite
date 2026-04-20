package drift

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// GroupedReport holds drift results organized by resource type.
type GroupedReport struct {
	Groups map[string]*TypeGroup `json:"groups"`
}

// TypeGroup holds resources of a single type.
type TypeGroup struct {
	Type     string     `json:"type"`
	Managed  []Resource `json:"managed"`
	Missing  []Resource `json:"missing"`
	Untracked []Resource `json:"untracked"`
}

// GroupByType organises a Report's resources by their type.
func GroupByType(r Report) GroupedReport {
	groups := make(map[string]*TypeGroup)

	ensure := func(t string) {
		if _, ok := groups[t]; !ok {
			groups[t] = &TypeGroup{Type: t}
		}
	}

	for _, res := range r.Managed {
		ensure(res.Type)
		groups[res.Type].Managed = append(groups[res.Type].Managed, res)
	}
	for _, res := range r.Missing {
		ensure(res.Type)
		groups[res.Type].Missing = append(groups[res.Type].Missing, res)
	}
	for _, res := range r.Untracked {
		ensure(res.Type)
		groups[res.Type].Untracked = append(groups[res.Type].Untracked, res)
	}

	return GroupedReport{Groups: groups}
}

// FprintGroup writes a human-readable grouped summary to w.
func FprintGroup(w io.Writer, g GroupedReport) {
	if len(g.Groups) == 0 {
		fmt.Fprintln(w, "No resources found.")
		return
	}

	// Sort types for deterministic output.
	types := make([]string, 0, len(g.Groups))
	for t := range g.Groups {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		grp := g.Groups[t]
		fmt.Fprintf(w, "[%s]\n", strings.ToUpper(t))
		fmt.Fprintf(w, "  managed:   %d\n", len(grp.Managed))
		if len(grp.Missing) > 0 {
			fmt.Fprintf(w, "  missing:   %d\n", len(grp.Missing))
			for _, res := range grp.Missing {
				fmt.Fprintf(w, "    - %s\n", res.ID)
			}
		}
		if len(grp.Untracked) > 0 {
			fmt.Fprintf(w, "  untracked: %d\n", len(grp.Untracked))
			for _, res := range grp.Untracked {
				fmt.Fprintf(w, "    + %s\n", res.ID)
			}
		}
	}
}
