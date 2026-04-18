package drift

import "strings"

// FilterOptions controls which resources are included in drift analysis.
type FilterOptions struct {
	ResourceTypes []string // if non-empty, only these types are scanned
	ExcludeTypes  []string // resource types to skip
}

// Filter returns a subset of resources matching the given options.
func Filter(resources []Resource, opts FilterOptions) []Resource {
	if len(opts.ResourceTypes) == 0 && len(opts.ExcludeTypes) == 0 {
		return resources
	}

	excludeSet := make(map[string]bool, len(opts.ExcludeTypes))
	for _, t := range opts.ExcludeTypes {
		excludeSet[strings.ToLower(t)] = true
	}

	includeSet := make(map[string]bool, len(opts.ResourceTypes))
	for _, t := range opts.ResourceTypes {
		includeSet[strings.ToLower(t)] = true
	}

	var out []Resource
	for _, r := range resources {
		rt := strings.ToLower(r.Type)
		if excludeSet[rt] {
			continue
		}
		if len(includeSet) > 0 && !includeSet[rt] {
			continue
		}
		out = append(out, r)
	}
	return out
}
