package drift

import (
	"fmt"
	"io"
	"sort"
)

// Annotation holds a user-defined note attached to a resource.
type Annotation struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Note    string `json:"note"`
	Author  string `json:"author,omitempty"`
	Created string `json:"created,omitempty"`
}

// AnnotationMap keys annotations by "type/id".
type AnnotationMap map[string]Annotation

// AnnotateReport attaches annotations to resources in a Report.
// Resources without a matching annotation are left unchanged.
func AnnotateReport(report Report, annotations AnnotationMap) Report {
	for i, r := range report.Missing {
		if a, ok := annotations[resourceKey(r)]; ok {
			if report.Missing[i].Attributes == nil {
				report.Missing[i].Attributes = map[string]string{}
			}
			report.Missing[i].Attributes["_annotation"] = a.Note
		}
	}
	for i, r := range report.Untracked {
		if a, ok := annotations[resourceKey(r)]; ok {
			if report.Untracked[i].Attributes == nil {
				report.Untracked[i].Attributes = map[string]string{}
			}
			report.Untracked[i].Attributes["_annotation"] = a.Note
		}
	}
	return report
}

// FprintAnnotations writes a human-readable list of annotations to w.
func FprintAnnotations(w io.Writer, annotations AnnotationMap) {
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No annotations.")
		return
	}
	keys := make([]string, 0, len(annotations))
	for k := range annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		a := annotations[k]
		line := fmt.Sprintf("  [%s] %s", k, a.Note)
		if a.Author != "" {
			line += fmt.Sprintf(" (by %s)", a.Author)
		}
		fmt.Fprintln(w, line)
	}
}

// AnnotationsHasEntries returns true when the map is non-empty.
func AnnotationsHasEntries(annotations AnnotationMap) bool {
	return len(annotations) > 0
}
