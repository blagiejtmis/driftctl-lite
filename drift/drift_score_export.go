package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// ExportScore writes the ScoreResult to path in the given format (json or csv).
func ExportScore(result ScoreResult, format, path string) error {
	switch format {
	case "json":
		return exportScoreJSON(result, path)
	case "csv":
		return exportScoreCSV(result, path)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func exportScoreJSON(result ScoreResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func exportScoreCSV(result ScoreResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	if err := w.Write([]string{"total", "managed", "missing", "untracked", "score_pct", "grade"}); err != nil {
		return err
	}
	return w.Write([]string{
		fmt.Sprintf("%d", result.Total),
		fmt.Sprintf("%d", result.Managed),
		fmt.Sprintf("%d", result.Missing),
		fmt.Sprintf("%d", result.Untracked),
		fmt.Sprintf("%.2f", result.ScorePct),
		result.Grade,
	})
}
