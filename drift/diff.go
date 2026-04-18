package drift

import (
	"fmt"
	"io"
	"sort"
)

// AttributeDiff holds a single attribute change between live and IaC state.
type AttributeDiff struct {
	Key      string
	IaC      string
	Live     string
}

// ResourceDiff holds all attribute differences for a resource.
type ResourceDiff struct {
	Resource   Resource
	Attributes []AttributeDiff
}

// Diff compares attribute-level differences between live and IaC resources
// that are present in both states (i.e. managed resources).
func Diff(live, iac []Resource) []ResourceDiff {
	liveIdx := indexResources(live)
	iacIdx := indexResources(iac)

	var diffs []ResourceDiff

	for key, iacRes := range iacIdx {
		liveRes, ok := liveIdx[key]
		if !ok {
			continue // missing resources handled by Compare
		}
		attrs := diffAttributes(liveRes.Attributes, iacRes.Attributes)
		if len(attrs) > 0 {
			diffs = append(diffs, ResourceDiff{
				Resource:   iacRes,
				Attributes: attrs,
			})
		}
	}

	sort.Slice(diffs, func(i, j int) bool {
		return resourceKey(diffs[i].Resource) < resourceKey(diffs[j].Resource)
	})
	return diffs
}

func diffAttributes(live, iac map[string]string) []AttributeDiff {
	var attrs []AttributeDiff
	keys := make(map[string]struct{})
	for k := range iac {
		keys[k] = struct{}{}
	}
	for k := range live {
		keys[k] = struct{}{}
	}
	for k := range keys {
		lv := live[k]
		iv := iac[k]
		if lv != iv {
			attrs = append(attrs, AttributeDiff{Key: k, IaC: iv, Live: lv})
		}
	}
	sort.Slice(attrs, func(i, j int) bool { return attrs[i].Key < attrs[j].Key })
	return attrs
}

// FprintDiff writes a human-readable attribute diff report to w.
func FprintDiff(w io.Writer, diffs []ResourceDiff) {
	if len(diffs) == 0 {
		fmt.Fprintln(w, "No attribute drift detected.")
		return
	}
	for _, d := range diffs {
		fmt.Fprintf(w, "[%s] %s\n", d.Resource.Type, d.Resource.ID)
		for _, a := range d.Attributes {
			fmt.Fprintf(w, "  ~ %s: IaC=%q  Live=%q\n", a.Key, a.IaC, a.Live)
		}
	}
}
