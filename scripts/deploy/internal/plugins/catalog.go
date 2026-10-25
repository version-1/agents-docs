package plugins

import (
	"deploy/internal/external"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Group struct {
	Skills      []string `json:"skills"`
	Includes    []string `json:"includes,omitempty"`
	Description string   `json:"description"`
}
type GroupsFile struct {
	Groups map[string]Group `json:"groups"`
}

type Catalog struct {
	Selected   map[string]string
	Individual []string
	Groups     []string
	Resolved   map[string][]string
}

var idPattern = regexp.MustCompile(`^[a-z0-9]+(?:[a-z0-9-]*[a-z0-9])?$`)

func load(path string, out any) error {
	b, e := os.ReadFile(path)
	if e != nil {
		return e
	}
	return json.Unmarshal(b, out)
}
func good(v string) bool { return idPattern.MatchString(v) }
func catalog(groups GroupsFile, externalPath string) (Catalog, error) {
	var zero Catalog
	selected := map[string]string{}
	names := make([]string, 0, len(groups.Groups))
	for n := range groups.Groups {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		if !good(n) {
			return zero, fmt.Errorf("invalid group %q", n)
		}
		if n == "individual-skills" || strings.HasPrefix(n, "skill-") {
			return zero, fmt.Errorf("reserved group %q", n)
		}
		for _, skill := range groups.Groups[n].Skills {
			if !good(skill) {
				return zero, fmt.Errorf("invalid skill %q", skill)
			}
			if old, ok := selected[skill]; ok {
				return zero, fmt.Errorf("skill %q is registered by groups %q and %q", skill, old, n)
			}
			selected[skill] = n
		}
		for _, included := range groups.Groups[n].Includes {
			if !good(included) {
				return zero, fmt.Errorf("invalid included group %q", included)
			}
			if _, ok := groups.Groups[included]; !ok {
				return zero, fmt.Errorf("unknown included group %q", included)
			}
		}
	}
	resolved, err := resolveGroups(groups, names)
	if err != nil {
		return zero, err
	}
	var locks []external.Skill
	if err := load(externalPath, &locks); err != nil {
		return zero, err
	}
	individual := make([]string, 0, len(locks))
	seen := map[string]bool{}
	for _, skill := range locks {
		if !good(skill.Name) || seen[skill.Name] {
			return zero, fmt.Errorf("invalid or duplicate external skill %q", skill.Name)
		}
		seen[skill.Name] = true
		individual = append(individual, skill.Name)
	}
	return Catalog{Selected: selected, Individual: individual, Groups: names, Resolved: resolved}, nil
}

func resolveGroups(groups GroupsFile, names []string) (map[string][]string, error) {
	resolved := map[string][]string{}
	visiting := map[string]bool{}
	var resolve func(string) ([]string, error)
	resolve = func(name string) ([]string, error) {
		if skills, ok := resolved[name]; ok {
			return skills, nil
		}
		if visiting[name] {
			return nil, fmt.Errorf("included group cycle at %q", name)
		}
		visiting[name] = true
		group := groups.Groups[name]
		skills := append([]string(nil), group.Skills...)
		for _, included := range group.Includes {
			includedSkills, err := resolve(included)
			if err != nil {
				return nil, err
			}
			skills = append(skills, includedSkills...)
		}
		visiting[name] = false
		sort.Strings(skills)
		resolved[name] = uniqueSorted(skills)
		return resolved[name], nil
	}
	for _, name := range names {
		if _, err := resolve(name); err != nil {
			return nil, err
		}
	}
	return resolved, nil
}
