package drift

import (
	"fmt"
	"io"
)

// ScoreResult holds the computed compliance score.
type ScoreResult struct {
	Score     float64
	Grade     string
	Managed   int
	Missing   int
	Untracked int
	Total     int
}

// ScoreReport computes a compliance score from a Report.
func ScoreReport(r Report) ScoreResult {
	total := len(r.Managed) + len(r.Missing) + len(r.Untracked)
	if total == 0 {
		return ScoreResult{Score: 100.0, Grade: "A", Total: 0}
	}
	score := float64(len(r.Managed)) / float64(total) * 100.0
	return ScoreResult{
		Score:     score,
		Grade:     grade(score),
		Managed:   len(r.Managed),
		Missing:   len(r.Missing),
		Untracked: len(r.Untracked),
		Total:     total,
	}
}

func grade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}

// FprintScore writes the score result to w.
func FprintScore(w io.Writer, s ScoreResult) {
	fmt.Fprintf(w, "Compliance Score: %.1f%% (Grade: %s)\n", s.Score, s.Grade)
	fmt.Fprintf(w, "  Managed:   %d\n", s.Managed)
	fmt.Fprintf(w, "  Missing:   %d\n", s.Missing)
	fmt.Fprintf(w, "  Untracked: %d\n", s.Untracked)
	fmt.Fprintf(w, "  Total:     %d\n", s.Total)
}
