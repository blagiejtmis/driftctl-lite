package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportForecast writes the ForecastResult to path in the given format (json or csv).
func ExportForecast(r ForecastResult, format, path string) error {
	switch strings.ToLower(format) {
	case "json":
		return exportForecastJSON(r, path)
	case "csv":
		return exportForecastCSV(r, path)
	default:
		return fmt.Errorf("unsupported forecast export format: %s", format)
	}
}

func exportForecastJSON(r ForecastResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("forecast export: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func exportForecastCSV(r ForecastResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("forecast export: %w", err)
	}
	defer f.Close()
	fmt.Fprintln(f, "period,score,grade,trend,confidence")
	for _, e := range r.Entries {
		fmt.Fprintf(f, "%s,%.2f,%s,%s,%s\n",
			csvEscape(e.Period), e.Score, csvEscape(e.Grade),
			csvEscape(r.Trend), csvEscape(r.Confidence))
	}
	return nil
}
