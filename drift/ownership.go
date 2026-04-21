package drift

import (
	"fmt"
	"io"
	"sort"
)

// OwnershipRule maps a resource type to a team owner.
type OwnershipRule struct {
	Type  string `json:"type"`
	Team  string `json:"team"`
	Email string `json:"email,omitempty"`
}

// OwnershipResult holds the resolved owner for a resource.
type OwnershipResult struct {
	Resource Resource
	Team     string
	Email    string
	Unowned  bool
}

// AssignOwnership maps each resource in the report to an owner based on rules.
// Resources with no matching rule are marked as unowned.
func AssignOwnership(report Report, rules []OwnershipRule) []OwnershipResult {
	index := make(map[string]OwnershipRule, len(rules))
	for _, r := range rules {
		index[r.Type] = r
	}

	all := append(append(append([]Resource{}, report.Managed...), report.Missing...), report.Untracked...)

	results := make([]OwnershipResult, 0, len(all))
	for _, res := range all {
		if rule, ok := index[res.Type]; ok {
			results = append(results, OwnershipResult{
				Resource: res,
				Team:     rule.Team,
				Email:    rule.Email,
			})
		} else {
			results = append(results, OwnershipResult{
				Resource: res,
				Unowned:  true,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Team != results[j].Team {
			return results[i].Team < results[j].Team
		}
		return results[i].Resource.ID < results[j].Resource.ID
	})

	return results
}

// OwnershipHasUnowned returns true if any resource has no assigned owner.
func OwnershipHasUnowned(results []OwnershipResult) bool {
	for _, r := range results {
		if r.Unowned {
			return true
		}
	}
	return false
}

// FprintOwnership writes a human-readable ownership report to w.
func FprintOwnership(w io.Writer, results []OwnershipResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No resources to report ownership for.")
		return
	}
	fmt.Fprintf(w, "%-30s %-20s %-20s %s\n", "RESOURCE", "TYPE", "TEAM", "EMAIL")
	for _, r := range results {
		team := r.Team
		if r.Unowned {
			team = "(unowned)"
		}
		fmt.Fprintf(w, "%-30s %-20s %-20s %s\n", r.Resource.ID, r.Resource.Type, team, r.Email)
	}
}
