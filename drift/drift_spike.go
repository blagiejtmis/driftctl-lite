package drift

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// SpikeEntry represents a detected drift spike for a resource type over a window.
type SpikeEntry struct {
	Type        string  `json:"type"`
	PrevDrift   int     `json:"prev_drift"`
	CurrDrift   int     `json:"curr_drift"`
	Delta       int     `json:"delta"`
	ChangePct   float64 `json:"change_pct"`
	IsSpike     bool    `json:"is_spike"`
}

// SpikeReport holds all spike evaluations.
type SpikeReport struct {
	Entries   []SpikeEntry `json:"entries"`
	Threshold float64      `json:"threshold"`
}

// SpikeHasSpikes returns true if any entry is flagged as a spike.
func SpikeHasSpikes(r SpikeReport) bool {
	for _, e := range r.Entries {
		if e.IsSpike {
			return true
		}
	}
	return false
}

// EvaluateSpike compares two velocity snapshots and flags types whose drift
// count grew by more than thresholdPct percent.
func EvaluateSpike(prev, curr map[string]int, thresholdPct float64) SpikeReport {
	types := make(map[string]struct{})
	for t := range prev {
		types[t] = struct{}{}
	}
	for t := range curr {
		types[t] = struct{}{}
	}

	var entries []SpikeEntry
	for t := range types {
		p := prev[t]
		c := curr[t]
		delta := c - p
		var pct float64
		if p > 0 {
			pct = float64(delta) / float64(p) * 100.0
		} else if c > 0 {
			pct = 100.0
		}
		entries = append(entries, SpikeEntry{
			Type:      t,
			PrevDrift: p,
			CurrDrift: c,
			Delta:     delta,
			ChangePct: pct,
			IsSpike:   pct >= thresholdPct && delta > 0,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ChangePct > entries[j].ChangePct
	})

	return SpikeReport{Entries: entries, Threshold: thresholdPct}
}

// FprintSpike writes a human-readable spike report to w.
func FprintSpike(w io.Writer, r SpikeReport) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "No drift spike data available.")
		return
	}
	fmt.Fprintf(w, "Drift Spike Report (threshold: %.0f%%)\n", r.Threshold)
	fmt.Fprintf(w, "%-30s %8s %8s %8s %10s %s\n", "TYPE", "PREV", "CURR", "DELTA", "CHANGE%", "SPIKE")
	for _, e := range r.Entries {
		spike := ""
		if e.IsSpike {
			spike = "⚠ SPIKE"
		}
		fmt.Fprintf(w, "%-30s %8d %8d %8d %9.1f%% %s\n",
			e.Type, e.PrevDrift, e.CurrDrift, e.Delta, e.ChangePct, spike)
	}
}

// FprintSpikeToStdout writes the spike report to stdout.
func FprintSpikeToStdout(r SpikeReport) {
	FprintSpike(os.Stdout, r)
}
