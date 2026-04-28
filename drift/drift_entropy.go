package drift

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
)

// EntropyEntry holds the drift entropy score for a resource type.
type EntropyEntry struct {
	Type       string  `json:"type"`
	Total      int     `json:"total"`
	Drifted    int     `json:"drifted"`
	Entropy    float64 `json:"entropy"`
	Normalized float64 `json:"normalized"` // 0.0 – 1.0
}

// EntropyReport is the result of BuildDriftEntropy.
type EntropyReport struct {
	Entries        []EntropyEntry `json:"entries"`
	OverallEntropy float64        `json:"overall_entropy"`
}

// BuildDriftEntropy computes Shannon entropy of drift distribution per type.
func BuildDriftEntropy(report Report) EntropyReport {
	type counts struct{ total, drifted int }
	buckets := map[string]*counts{}

	for _, r := range report.Managed {
		if buckets[r.Type] == nil {
			buckets[r.Type] = &counts{}
		}
		buckets[r.Type].total++
	}
	for _, r := range report.Missing {
		if buckets[r.Type] == nil {
			buckets[r.Type] = &counts{}
		}
		buckets[r.Type].total++
		buckets[r.Type].drifted++
	}
	for _, r := range report.Untracked {
		if buckets[r.Type] == nil {
			buckets[r.Type] = &counts{}
		}
		buckets[r.Type].total++
		buckets[r.Type].drifted++
	}

	entries := make([]EntropyEntry, 0, len(buckets))
	for typ, c := range buckets {
		e := entropy(c.drifted, c.total)
		var norm float64
		if c.total > 0 {
			norm = float64(c.drifted) / float64(c.total)
		}
		entries = append(entries, EntropyEntry{
			Type:       typ,
			Total:      c.total,
			Drifted:    c.drifted,
			Entropy:    e,
			Normalized: norm,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Entropy != entries[j].Entropy {
			return entries[i].Entropy > entries[j].Entropy
		}
		return entries[i].Type < entries[j].Type
	})

	overall := 0.0
	for _, e := range entries {
		overall += e.Entropy
	}

	return EntropyReport{Entries: entries, OverallEntropy: overall}
}

// entropy computes binary Shannon entropy for k drifted out of n total.
func entropy(k, n int) float64 {
	if n == 0 || k == 0 || k == n {
		return 0.0
	}
	p := float64(k) / float64(n)
	q := 1.0 - p
	return -(p*math.Log2(p) + q*math.Log2(q))
}

// EntropyHasDrift returns true when any type has drifted resources.
func EntropyHasDrift(r EntropyReport) bool {
	for _, e := range r.Entries {
		if e.Drifted > 0 {
			return true
		}
	}
	return false
}

// FprintEntropy writes a human-readable entropy report to w.
func FprintEntropy(w io.Writer, r EntropyReport) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "[entropy] no resources evaluated")
		return
	}
	fmt.Fprintf(w, "[entropy] overall=%.4f\n", r.OverallEntropy)
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-30s total=%-4d drifted=%-4d entropy=%.4f norm=%.2f\n",
			e.Type, e.Total, e.Drifted, e.Entropy, e.Normalized)
	}
}

// FprintEntropyToStdout is a convenience wrapper.
func FprintEntropyToStdout(r EntropyReport) {
	FprintEntropy(os.Stdout, r)
}
