package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportMaturity writes the maturity report to path in the given format (json|csv).
func ExportMaturity(r MaturityReport, format, path string) error {
	switch strings.ToLower(format) {
	case "json":
		return exportMaturityJSON(r, path)
	case "csv":
		return exportMaturityCSV(r, path)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func exportMaturityJSON(r MaturityReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func exportMaturityCSV(r MaturityReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, "type,total,managed,coverage_pct,level")
	for _, res := range r.Results {
		fmt.Fprintf(f, "%s,%d,%d,%.2f,%s\n",
			csvEscape(res.Type), res.Total, res.Managed, res.Coverage, res.Level)
	}
	return nil
}
