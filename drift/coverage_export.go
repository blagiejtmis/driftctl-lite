package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportCoverageOptions configures coverage export.
type ExportCoverageOptions struct {
	Format string
	Path   string
}

// ExportCoverage writes the coverage report to a file in JSON or CSV format.
func ExportCoverage(cr CoverageReport, opts ExportCoverageOptions) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return exportCoverageJSON(cr, opts.Path)
	case "csv":
		return exportCoverageCSV(cr, opts.Path)
	default:
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

func exportCoverageJSON(cr CoverageReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cr)
}

func exportCoverageCSV(cr CoverageReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "type,total,managed,coverage_pct")
	for _, r := range cr.Results {
		fmt.Fprintf(f, "%s,%d,%d,%.2f\n",
			csvEscape(r.Type), r.Total, r.Managed, r.Pct)
	}
	return nil
}
