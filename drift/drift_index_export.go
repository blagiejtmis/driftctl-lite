package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportDriftIndex writes the DriftIndex to path in the given format (json|csv).
func ExportDriftIndex(idx DriftIndex, format, path string) error {
	switch strings.ToLower(format) {
	case "json":
		return exportDriftIndexJSON(idx, path)
	case "csv":
		return exportDriftIndexCSV(idx, path)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func exportDriftIndexJSON(idx DriftIndex, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(idx)
}

func exportDriftIndexCSV(idx DriftIndex, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, "type,total,managed,missing,untracked,changed,drift_pct")
	for _, e := range idx.Entries {
		fmt.Fprintf(f, "%s,%d,%d,%d,%d,%d,%.2f\n",
			csvEscape(e.Type), e.Total, e.Managed, e.Missing, e.Untracked, e.Changed, e.DriftPct)
	}
	// aggregate row
	agg := idx.Total
	fmt.Fprintf(f, "%s,%d,%d,%d,%d,%d,%.2f\n",
		csvEscape(agg.Type), agg.Total, agg.Managed, agg.Missing, agg.Untracked, agg.Changed, agg.DriftPct)
	return nil
}
