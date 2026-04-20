package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// ExportRisk writes risk results to a file in the given format (json or csv).
func ExportRisk(results []RiskResult, format, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export risk: create file: %w", err)
	}
	defer f.Close()

	switch format {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(results); err != nil {
			return fmt.Errorf("export risk: encode json: %w", err)
		}
	case "csv":
		w := csv.NewWriter(f)
		_ = w.Write([]string{"resource_id", "resource_type", "level", "reason", "score"})
		for _, r := range results {
			_ = w.Write([]string{
				r.ResourceID,
				r.ResourceType,
				r.Level,
				r.Reason,
				strconv.Itoa(r.Score),
			})
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return fmt.Errorf("export risk: write csv: %w", err)
		}
	default:
		return fmt.Errorf("export risk: unsupported format %q", format)
	}
	return nil
}
