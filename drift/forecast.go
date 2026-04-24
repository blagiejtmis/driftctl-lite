package drift

import (
	"fmt"
	"io"
	"math"
	"sort"
)

// ForecastEntry represents a predicted drift score at a future point.
type ForecastEntry struct {
	Period string  `json:"period"`
	Score  float64 `json:"score"`
	Grade  string  `json:"grade"`
}

// ForecastResult holds the full forecast output.
type ForecastResult struct {
	Entries    []ForecastEntry `json:"entries"`
	Trend      string          `json:"trend"`
	Confidence string          `json:"confidence"`
}

// Forecast predicts future drift scores based on historical trend entries.
// It uses simple linear regression over the last N scores.
func Forecast(history []TrendEntry, periods int) ForecastResult {
	if len(history) == 0 || periods <= 0 {
		return ForecastResult{Trend: "unknown", Confidence: "none"}
	}

	n := len(history)
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i, e := range history {
		xs[i] = float64(i)
		ys[i] = e.Score
	}

	slope, intercept := linearRegression(xs, ys)

	result := ForecastResult{}
	for p := 1; p <= periods; p++ {
		predicted := intercept + slope*float64(n-1+p)
		predicted = math.Max(0, math.Min(100, predicted))
		result.Entries = append(result.Entries, ForecastEntry{
			Period: fmt.Sprintf("+%d", p),
			Score:  math.Round(predicted*100) / 100,
			Grade:  grade(predicted),
		})
	}

	result.Trend = trendDirection(slope)
	result.Confidence = confidenceLevel(n)
	return result
}

func linearRegression(xs, ys []float64) (slope, intercept float64) {
	n := float64(len(xs))
	var sumX, sumY, sumXY, sumXX float64
	for i := range xs {
		sumX += xs[i]
		sumY += ys[i]
		sumXY += xs[i] * ys[i]
		sumXX += xs[i] * xs[i]
	}
	denom := n*sumXX - sumX*sumX
	if denom == 0 {
		return 0, sumY / n
	}
	slope = (n*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / n
	return
}

func trendDirection(slope float64) string {
	switch {
	case slope > 0.5:
		return "improving"
	case slope < -0.5:
		return "degrading"
	default:
		return "stable"
	}
}

func confidenceLevel(n int) string {
	switch {
	case n >= 10:
		return "high"
	case n >= 5:
		return "medium"
	default:
		return "low"
	}
}

// ForecastHasDrift returns true if any forecast period predicts score below 80.
func ForecastHasDrift(r ForecastResult) bool {
	for _, e := range r.Entries {
		if e.Score < 80 {
			return true
		}
	}
	return false
}

// FprintForecast writes a human-readable forecast table to w.
func FprintForecast(w io.Writer, r ForecastResult) {
	fmt.Fprintf(w, "Drift Forecast (trend: %s, confidence: %s)\n", r.Trend, r.Confidence)
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "  No forecast available.")
		return
	}
	fmt.Fprintf(w, "  %-10s %-8s %s\n", "Period", "Score", "Grade")
	sort.Slice(r.Entries, func(i, j int) bool {
		return r.Entries[i].Period < r.Entries[j].Period
	})
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-10s %-8.2f %s\n", e.Period, e.Score, e.Grade)
	}
}
