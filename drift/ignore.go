package drift

import (
	"bufio"
	"os"
	"strings"
)

// IgnoreRule represents a single ignore entry (type + optional id glob).
type IgnoreRule struct {
	Type string
	ID   string // empty means match all IDs of that type
}

// LoadIgnoreFile parses a .driftignore file and returns a list of rules.
// Each non-blank, non-comment line must be: "type" or "type/id".
func LoadIgnoreFile(path string) ([]IgnoreRule, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var rules []IgnoreRule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "/", 2)
		rule := IgnoreRule{Type: parts[0]}
		if len(parts) == 2 {
			rule.ID = parts[1]
		}
		rules = append(rules, rule)
	}
	return rules, scanner.Err()
}

// ApplyIgnore removes resources that match any ignore rule.
func ApplyIgnore(resources []Resource, rules []IgnoreRule) []Resource {
	if len(rules) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if !isIgnored(r, rules) {
			out = append(out, r)
		}
	}
	return out
}

func isIgnored(r Resource, rules []IgnoreRule) bool {
	for _, rule := range rules {
		if !strings.EqualFold(rule.Type, r.Type) {
			continue
		}
		if rule.ID == "" || rule.ID == r.ID {
			return true
		}
	}
	return false
}
