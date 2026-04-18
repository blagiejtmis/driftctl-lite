package drift

import "fmt"

// DriftScore represents a numeric health score for the scanned infrastructure.
type DriftScore struct {
	Total     int     `json:"total"`
	Managed   int     `json:"managed"`
	Missing   int     `json:"missing"`
	Untracked int     `json:"untracked"`
	Score     float64 `json:"score"` // 0.0 (all drift) to 100.0 (no drift)
	Grade     string  `json:"grade"`
}

// ScoreReport computes a drift health score from a Report.
func ScoreReport(r Report) DriftScore {
	total := len(r.Managed) + len(r.Missing) + len(r.Untracked)
	if total == 0 {
		return DriftScore{Total: 0, Score: 100.0, Grade: "A"}
	}
	managed := len(r.Managed)
	score := (float64(managed) / float64(total)) * 100.0
	return DriftScore{
		Total:     total,
		Managed:   managed,
		Missing:   len(r.Missing),
		Untracked: len(r.Untracked),
		Score:     score,
		Grade:     grade(score),
	}
}

func grade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 50:
		return "C"
	case score >= 25:
		return "D"
	default:
		return "F"
	}
}

// FprintScore writes a human-readable score summary to w.
func FprintScore(w interface{ WriteString(string) (int, error) }, s DriftScore) {
	w.WriteString(fmt.Sprintf("Drift Score: %.1f / 100  (Grade: %s)\n", s.Score, s.Grade))
	w.WriteString(fmt.Sprintf("  Total resources : %d\n", s.Total))
	w.WriteString(fmt.Sprintf("  Managed         : %d\n", s.Managed))
	w.WriteString(fmt.Sprintf("  Missing         : %d\n", s.Missing))
	w.WriteString(fmt.Sprintf("  Untracked       : %d\n", s.Untracked))
}
