package drift

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// DepEdge represents a directional dependency between two resource keys.
type DepEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// DepGraph holds adjacency information for resource dependencies.
type DepGraph struct {
	Edges []DepEdge          `json:"edges"`
	adj   map[string][]string
}

// BuildDepGraph constructs a dependency graph from a list of resources.
// Resources may declare dependencies via the "depends_on" attribute as a
// comma-separated list of resource keys.
func BuildDepGraph(resources []Resource) *DepGraph {
	g := &DepGraph{adj: make(map[string][]string)}
	for _, r := range resources {
		key := resourceKey(r)
		raw, ok := r.Attributes["depends_on"]
		if !ok {
			continue
		}
		for _, dep := range strings.Split(fmt.Sprintf("%v", raw), ",") {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				continue
			}
			g.Edges = append(g.Edges, DepEdge{From: key, To: dep})
			g.adj[key] = append(g.adj[key], dep)
		}
	}
	return g
}

// Affected returns all resource keys transitively affected when the given key
// changes (i.e. dependents that depend on it, resolved via reverse traversal).
func (g *DepGraph) Affected(key string) []string {
	// Build reverse adjacency map.
	rev := make(map[string][]string)
	for _, e := range g.Edges {
		rev[e.To] = append(rev[e.To], e.From)
	}
	visited := make(map[string]bool)
	var walk func(k string)
	walk = func(k string) {
		for _, parent := range rev[k] {
			if !visited[parent] {
				visited[parent] = true
				walk(parent)
			}
		}
	}
	walk(key)
	result := make([]string, 0, len(visited))
	for k := range visited {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

// FprintDeps writes a human-readable dependency graph to w.
func FprintDeps(w io.Writer, g *DepGraph) {
	if len(g.Edges) == 0 {
		fmt.Fprintln(w, "No resource dependencies found.")
		return
	}
	fmt.Fprintln(w, "Resource Dependency Graph:")
	for _, e := range g.Edges {
		fmt.Fprintf(w, "  %s  -->  %s\n", e.From, e.To)
	}
}
