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
	Description string   `json:"description"`
}
type GroupsFile struct {
	Groups map[string]Group `json:"groups"`
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
func allSkills(groups GroupsFile, externalPath string) (map[string]string, []string, []string, error) {
	selected := map[string]string{}
	names := make([]string, 0, len(groups.Groups))
	for n := range groups.Groups {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		if !good(n) {
			return nil, nil, nil, fmt.Errorf("invalid group %q", n)
		}
		if n == "individual-skills" || strings.HasPrefix(n, "skill-") {
			return nil, nil, nil, fmt.Errorf("reserved group %q", n)
		}
		for _, skill := range groups.Groups[n].Skills {
			if !good(skill) {
				return nil, nil, nil, fmt.Errorf("invalid skill %q", skill)
			}
			if old, ok := selected[skill]; ok {
				return nil, nil, nil, fmt.Errorf("skill %q is registered by groups %q and %q", skill, old, n)
			}
			selected[skill] = n
		}
	}
	var locks []external.Skill
	if err := load(externalPath, &locks); err != nil {
		return nil, nil, nil, err
	}
	individual := make([]string, 0, len(locks))
	seen := map[string]bool{}
	for _, skill := range locks {
		if !good(skill.Name) || seen[skill.Name] {
			return nil, nil, nil, fmt.Errorf("invalid or duplicate external skill %q", skill.Name)
		}
		seen[skill.Name] = true
		individual = append(individual, skill.Name)
	}
	return selected, individual, names, nil
}
