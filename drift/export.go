package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat represents supported export formats.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)

// Export writes the drift report to a file in the given format.
func Export(report Report, path string, format ExportFormat) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: cannot create file %q: %w", path, err)
	}
	defer f.Close()

	switch format {
	case ExportJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return fmt.Errorf("export: json encode: %w", err)
		}
	case ExportCSV:
		fmt.Fprintln(f, "status,type,id")
		for _, r := range report.Missing {
			fmt.Fprintf(f, "missing,%s,%s\n", csvEscape(r.Type), csvEscape(r.ID))
		}
		for _, r := range report.Untracked {
			fmt.Fprintf(f, "untracked,%s,%s\n", csvEscape(r.Type), csvEscape(r.ID))
		}
		for _, r := range report.Managed {
			fmt.Fprintf(f, "managed,%s,%s\n", csvEscape(r.Type), csvEscape(r.ID))
		}
	default:
		return fmt.Errorf("export: unsupported format %q", format)
	}
	return nil
}

func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + s + "\""
	}
	return s
}
