package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// ExportSurface writes the SurfaceReport to path in the given format (json|csv).
func ExportSurface(r SurfaceReport, format, path string) error {
	switch format {
	case "json":
		return exportSurfaceJSON(r, path)
	case "csv":
		return exportSurfaceCSV(r, path)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func exportSurfaceJSON(r SurfaceReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func exportSurfaceCSV(r SurfaceReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	_ = w.Write([]string{"type", "total", "drifted", "surface_pct", "risk_level"})
	for _, e := range r.Entries {
		_ = w.Write([]string{
			e.Type,
			strconv.Itoa(e.Total),
			strconv.Itoa(e.Drifted),
			strconv.FormatFloat(e.SurfacePct, 'f', 2, 64),
			e.RiskLevel,
		})
	}
	w.Flush()
	return w.Error()
}
