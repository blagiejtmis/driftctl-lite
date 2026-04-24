package drift

import (
	"bytes"
	"testing"
)

func baseForecastHistory() []TrendEntry {
	return []TrendEntry{
		{Score: 60},
		{Score: 65},
		{Score: 70},
		{Score: 75},
		{Score: 80},
	}
}

func TestForecast_EmptyHistory(t *testing.T) {
	r := Forecast(nil, 3)
	if r.Trend != "unknown" {
		t.Errorf("expected unknown trend, got %s", r.Trend)
	}
	if len(r.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(r.Entries))
	}
}

func TestForecast_ZeroPeriods(t *testing.T) {
	r := Forecast(baseForecastHistory(), 0)
	if len(r.Entries) != 0 {
		t.Errorf("expected no entries for 0 periods")
	}
}

func TestForecast_CorrectPeriodCount(t *testing.T) {
	r := Forecast(baseForecastHistory(), 3)
	if len(r.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(r.Entries))
	}
}

func TestForecast_ImprovingTrend(t *testing.T) {
	r := Forecast(baseForecastHistory(), 2)
	if r.Trend != "improving" {
		t.Errorf("expected improving trend, got %s", r.Trend)
	}
}

func TestForecast_DegradingTrend(t *testing.T) {
	h := []TrendEntry{{Score: 90}, {Score: 80}, {Score: 70}, {Score: 60}, {Score: 50}}
	r := Forecast(h, 2)
	if r.Trend != "degrading" {
		t.Errorf("expected degrading trend, got %s", r.Trend)
	}
}

func TestForecast_ConfidenceLow(t *testing.T) {
	r := Forecast(baseForecastHistory(), 1)
	if r.Confidence != "medium" {
		t.Errorf("expected medium confidence for 5 points, got %s", r.Confidence)
	}
}

func TestForecastHasDrift_True(t *testing.T) {
	r := ForecastResult{Entries: []ForecastEntry{{Score: 75, Grade: "C"}}}
	if !ForecastHasDrift(r) {
		t.Error("expected drift")
	}
}

func TestForecastHasDrift_False(t *testing.T) {
	r := ForecastResult{Entries: []ForecastEntry{{Score: 95, Grade: "A"}}}
	if ForecastHasDrift(r) {
		t.Error("expected no drift")
	}
}

func TestFprintForecast_Empty(t *testing.T) {
	var buf bytes.Buffer
	FprintForecast(&buf, ForecastResult{Trend: "unknown", Confidence: "none"})
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFprintForecast_WithEntries(t *testing.T) {
	r := Forecast(baseForecastHistory(), 2)
	var buf bytes.Buffer
	FprintForecast(&buf, r)
	out := buf.String()
	if len(out) < 10 {
		t.Error("expected meaningful output")
	}
}
