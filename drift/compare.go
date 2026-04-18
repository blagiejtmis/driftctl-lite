package drift

// Compare takes a slice of IaC resources and a slice of live resources,
// returning a DriftResult describing the differences.
func Compare(iac []Resource, live []Resource) DriftResult {
	result := DriftResult{}

	liveIndex := indexResources(live)
	iacIndex := indexResources(iac)

	// Check each IaC resource against live state.
	for _, r := range iac {
		key := resourceKey(r)
		if _, found := liveIndex[key]; found {
			result.Managed = append(result.Managed, r)
		} else {
			result.Missing = append(result.Missing, r)
		}
	}

	// Find live resources not tracked in IaC.
	for _, r := range live {
		key := resourceKey(r)
		if _, found := iacIndex[key]; !found {
			result.Untracked = append(result.Untracked, r)
		}
	}

	return result
}

func indexResources(resources []Resource) map[string]Resource {
	idx := make(map[string]Resource, len(resources))
	for _, r := range resources {
		idx[resourceKey(r)] = r
	}
	return idx
}

func resourceKey(r Resource) string {
	return string(r.Type) + "/" + r.ID
}
