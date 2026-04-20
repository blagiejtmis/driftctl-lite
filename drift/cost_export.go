package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// ExportCost writes a CostReport to path in the given format (json or csv).
func ExportCost(cr CostReport, format, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cost export: create file: %w", err)
	}
	defer f.Close()

	switch format {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(cr); err != nil {
			return fmt.Errorf("cost export: encode json: %w", err)
		}
	case "csv":
		w := csv.NewWriter(f)
		_ = w.Write([]string{"resource_id", "resource_type", "estimated_usd", "drifted"})
		for _, e := range cr.Entries {
			_ = w.Write([]string{
				e.ResourceID,
				e.ResourceType,
				strconv.FormatFloat(e.EstimatedUSD, 'f', 2, 64),
				strconv.FormatBool(e.Drifted),
			})
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return fmt.Errorf("cost export: write csv: %w", err)
		}
	default:
		return fmt.Errorf("cost export: unsupported format %q", format)
	}
	return nil
}
