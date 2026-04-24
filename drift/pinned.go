package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// PinnedEntry records a resource whose current state has been pinned (frozen).
type PinnedEntry struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

// PinnedSet is the full set of pinned resources.
type PinnedSet struct {
	Entries []PinnedEntry `json:"pinned"`
}

// PinResources builds a PinnedSet from the provided resources.
func PinResources(resources []Resource) PinnedSet {
	pinned := make([]PinnedEntry, 0, len(resources))
	for _, r := range resources {
		attrs := make(map[string]string, len(r.Attributes))
		for k, v := range r.Attributes {
			attrs[k] = v
		}
		pinned = append(pinned, PinnedEntry{
			Type:       r.Type,
			ID:         r.ID,
			Attributes: attrs,
		})
	}
	sort.Slice(pinned, func(i, j int) bool {
		if pinned[i].Type != pinned[j].Type {
			return pinned[i].Type < pinned[j].Type
		}
		return pinned[i].ID < pinned[j].ID
	})
	return PinnedSet{Entries: pinned}
}

// SavePinned writes a PinnedSet to the given file path as JSON.
func SavePinned(path string, ps PinnedSet) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("pinned: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(ps)
}

// LoadPinned reads a PinnedSet from the given file path.
func LoadPinned(path string) (PinnedSet, error) {
	f, err := os.Open(path)
	if err != nil {
		return PinnedSet{}, fmt.Errorf("pinned: open %s: %w", path, err)
	}
	defer f.Close()
	var ps PinnedSet
	if err := json.NewDecoder(f).Decode(&ps); err != nil {
		return PinnedSet{}, fmt.Errorf("pinned: decode %s: %w", path, err)
	}
	return ps, nil
}

// FprintPinned writes a human-readable summary of the pinned set to w.
func FprintPinned(w io.Writer, ps PinnedSet) {
	if len(ps.Entries) == 0 {
		fmt.Fprintln(w, "No pinned resources.")
		return
	}
	fmt.Fprintf(w, "Pinned resources (%d):\n", len(ps.Entries))
	for _, e := range ps.Entries {
		fmt.Fprintf(w, "  [%s] %s\n", e.Type, e.ID)
	}
}
