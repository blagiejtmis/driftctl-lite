package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// OutputFormat defines the output format for scan results.
type OutputFormat string

const (
	OutputText OutputFormat = "text"
	OutputJSON OutputFormat = "json"
)

// Report holds the result of a drift comparison.
type Report struct {
	Managed   []Resource `json:"managed"`
	Missing   []Resource `json:"missing"`
	Untracked []Resource `json:"untracked"`
}

// HasDrift returns true if any drift was detected.
func (r *Report) HasDrift() bool {
	return len(r.Missing) > 0 || len(r.Untracked) > 0
}

// Print writes the report to stdout in the requested format.
func (r *Report) Print(format OutputFormat) error {
	return r.Fprint(os.Stdout, format)
}

// Fprint writes the report to the given writer.
func (r *Report) Fprint(w io.Writer, format OutputFormat) error {
	switch format {
	case OutputJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	default:
		return r.fprintText(w)
	}
}

func (r *Report) fprintText(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "STATUS\tTYPE\tID\n")
	fmt.Fprintf(tw, "------\t----\t--\n")
	for _, res := range r.Managed {
		fmt.Fprintf(tw, "OK\t%s\t%s\n", res.Type, res.ID)
	}
	for _, res := range r.Missing {
		fmt.Fprintf(tw, "MISSING\t%s\t%s\n", res.Type, res.ID)
	}
	for _, res := range r.Untracked {
		fmt.Fprintf(tw, "UNTRACKED\t%s\t%s\n", res.Type, res.ID)
	}
	return tw.Flush()
}
