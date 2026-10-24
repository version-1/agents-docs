package plugins

import (
	"crypto/sha256"
	"deploy/internal/external"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Options struct {
	GroupsPath, ExternalPath, Output string
	Fetcher                          external.Fetcher
	DryRun                           bool
}
type Summary struct {
	Groups  []string `json:"groups"`
	Skills  []string `json:"skills"`
	Targets []string `json:"targets"`
}

var siblingReference = regexp.MustCompile(`\.\./([a-z0-9][a-z0-9-]*)/[^\s)"']+`)

func has(v []string, x string) bool {
	for _, y := range v {
		if y == x {
			return true
		}
	}
	return false
}

func Build(o Options) (Summary, error) {
	var zero Summary
	var gs GroupsFile
	if err := load(o.GroupsPath, &gs); err != nil {
		return zero, fmt.Errorf("load groups: %w", err)
	}
	selected, individual, groups, err := allSkills(gs, o.ExternalPath)
	if err != nil {
		return zero, err
	}
	targets := []string{"codex", "claude"}
	if o.ExternalPath == "" {
		return zero, fmt.Errorf("external skills is required")
	}
	if len(selected) == 0 && len(individual) == 0 {
		return zero, fmt.Errorf("selection is empty")
	}
	sources, cleanup, err := sources(filepath.Join(filepath.Dir(o.GroupsPath), "skills", "internal"), o.ExternalPath, selected, individual, o.Fetcher)
	if err != nil {
		return zero, err
	}
	defer cleanup()
	root, e := filepath.Abs(o.Output)
	if e != nil {
		return zero, e
	}
	if e = os.MkdirAll(filepath.Dir(root), 0755); e != nil {
		return zero, e
	}
	if info, statErr := os.Lstat(root); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		return zero, fmt.Errorf("refuse symlink output %q", root)
	}
	stage, e := os.MkdirTemp(filepath.Dir(root), ".marketplace-stage-")
	if e != nil {
		return zero, e
	}
	defer os.RemoveAll(stage)
	packs := map[string][]string{}
	for skill, g := range selected {
		packs[g] = append(packs[g], skill)
		packs["skill-"+skill] = []string{skill}
	}
	if len(individual) > 0 {
		for _, skill := range individual {
			packs["skill-"+skill] = []string{skill}
		}
	}
	for p, skills := range packs {
		sort.Strings(skills)
		if e = writePlugin(stage, p, skills, sources); e != nil {
			return zero, e
		}
	}
	sum := Summary{Groups: groups, Targets: targets}
	for skill := range selected {
		sum.Skills = append(sum.Skills, skill)
	}
	sum.Skills = append(sum.Skills, individual...)
	sort.Strings(sum.Skills)
	sum.Skills = uniqueSorted(sum.Skills)
	if has(targets, "codex") {
		e = codex(stage, packs)
	}
	if e == nil && has(targets, "claude") {
		e = claude(stage, packs)
	}
	if e == nil {
		e = write(filepath.Join(stage, "summary.json"), sum)
	}
	if e != nil {
		return zero, e
	}
	if o.DryRun {
		return sum, nil
	}
	if _, e = os.Stat(root); e == nil {
		if _, e = os.Stat(filepath.Join(root, ".agents-docs-marketplace")); e != nil {
			return zero, fmt.Errorf("refuse to replace unmarked output %q", root)
		}
	}
	if e = os.WriteFile(filepath.Join(stage, ".agents-docs-marketplace"), []byte("generated\n"), 0644); e != nil {
		return zero, e
	}
	backup := fmt.Sprintf("%s.previous-%d", root, time.Now().UnixNano())
	if _, e = os.Stat(root); e == nil {
		if e = os.Rename(root, backup); e != nil {
			return zero, e
		}
	}
	if e = os.Rename(stage, root); e != nil {
		_ = os.Rename(backup, root)
		return zero, e
	}
	_ = os.RemoveAll(backup)
	return sum, nil
}

func uniqueSorted(values []string) []string {
	if len(values) == 0 {
		return values
	}
	result := values[:1]
	for _, value := range values[1:] {
		if value != result[len(result)-1] {
			result = append(result, value)
		}
	}
	return result
}
func sources(internal, externalPath string, selected map[string]string, individual []string, fetcher external.Fetcher) (map[string]string, func(), error) {
	wanted := map[string]bool{}
	for s := range selected {
		wanted[s] = true
	}
	for _, s := range individual {
		wanted[s] = true
	}
	var ext []external.Skill
	var err error
	if externalPath != "" {
		ext, err = external.Load(externalPath)
		if err != nil {
			return nil, func() {}, err
		}
	}
	byName := map[string]external.Skill{}
	for _, s := range ext {
		if _, exists := byName[s.Name]; exists {
			cleanup := func() {}
			return nil, cleanup, fmt.Errorf("duplicate external skill %q", s.Name)
		}
		byName[s.Name] = s
	}
	tmp, err := os.MkdirTemp("", "plugin-skills-")
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() { _ = os.RemoveAll(tmp) }
	result := map[string]string{}
	for name := range wanted {
		p := filepath.Join(internal, name)
		info, e := os.Lstat(p)
		if e == nil {
			if _, externalExists := byName[name]; externalExists {
				cleanup()
				return nil, func() {}, fmt.Errorf("internal and external skill name conflict %q", name)
			}
			if info.Mode()&os.ModeSymlink != 0 {
				cleanup()
				return nil, func() {}, fmt.Errorf("skill %q is a symlink", name)
			}
			if e = external.ValidateFetched(external.Skill{Name: name}, p); e != nil {
				cleanup()
				return nil, func() {}, e
			}
			result[name] = p
			continue
		}
		skill, ok := byName[name]
		if !ok {
			cleanup()
			return nil, func() {}, fmt.Errorf("unknown skill %q", name)
		}
		if fetcher == nil {
			fetcher = external.GitFetcher{}
		}
		p, e = fetcher.Fetch(skill, tmp)
		if e != nil {
			cleanup()
			return nil, func() {}, e
		}
		if e = external.ValidateFetched(skill, p); e != nil {
			cleanup()
			return nil, func() {}, e
		}
		result[name] = p
	}
	return result, cleanup, nil
}
func writePlugin(root, plugin string, skills []string, src map[string]string) error {
	if strings.HasPrefix(plugin, "skill-") {
		var err error
		skills, err = skillClosure(skills, src)
		if err != nil {
			return err
		}
	}
	for _, s := range skills {
		if e := copyTree(src[s], filepath.Join(root, "plugins", plugin, "skills", s)); e != nil {
			return e
		}
	}
	version, e := pluginVersion(skills, src)
	if e != nil {
		return e
	}
	if e := write(filepath.Join(root, "plugins", plugin, ".codex-plugin", "plugin.json"), map[string]any{"name": plugin, "version": version, "skills": "./skills/"}); e != nil {
		return e
	}
	return write(filepath.Join(root, "plugins", plugin, ".claude-plugin", "plugin.json"), map[string]any{"name": plugin, "version": version, "description": plugin})
}

func pluginVersion(skills []string, sources map[string]string) (string, error) {
	h := sha256.New()
	for _, skill := range skills {
		if _, err := h.Write([]byte(skill)); err != nil {
			return "", err
		}
		err := filepath.WalkDir(sources[skill], func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(sources[skill], path)
			if err != nil {
				return err
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			if _, err = h.Write([]byte(rel + info.Mode().Perm().String())); err != nil {
				return err
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = h.Write(b)
			return err
		})
		if err != nil {
			return "", err
		}
	}
	return "0.0.0-" + fmt.Sprintf("%x", h.Sum(nil))[:12], nil
}

// skillClosure follows only sibling skill references such as ../code-general/SKILL.md.
func skillClosure(initial []string, sources map[string]string) ([]string, error) {
	seen := map[string]bool{}
	queue := append([]string(nil), initial...)
	for len(queue) > 0 {
		skill := queue[0]
		queue = queue[1:]
		if seen[skill] {
			continue
		}
		seen[skill] = true
		root, ok := sources[skill]
		if !ok {
			return nil, fmt.Errorf("referenced skill %q is not available", skill)
		}
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, match := range siblingReference.FindAllStringSubmatch(string(b), -1) {
				if _, exists := sources[match[1]]; exists && !seen[match[1]] {
					queue = append(queue, match[1])
				} else if strings.HasSuffix(match[0], "/SKILL.md") && strings.HasPrefix(match[1], "skill-") {
					return fmt.Errorf("referenced skill %q is not available", match[1])
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	result := make([]string, 0, len(seen))
	for skill := range seen {
		result = append(result, skill)
	}
	sort.Strings(result)
	return result, nil
}
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("refuse symlink %q", p)
		}
		rel, e := filepath.Rel(src, p)
		if e != nil {
			return e
		}
		to := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(to, 0755)
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		if e = os.MkdirAll(filepath.Dir(to), 0755); e != nil {
			return e
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		return os.WriteFile(to, b, info.Mode().Perm())
	})
}
func write(path string, v any) error {
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(path), 0755); e != nil {
		return e
	}
	return os.WriteFile(path, append(b, '\n'), 0644)
}
func codex(root string, packs map[string][]string) error {
	list := []any{}
	names := make([]string, 0, len(packs))
	for n := range packs {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		list = append(list, map[string]any{"name": n, "source": map[string]string{"source": "local", "path": "./plugins/" + n}, "policy": map[string]string{"installation": "AVAILABLE", "authentication": "ON_INSTALL"}, "category": "Productivity"})
	}
	return write(filepath.Join(root, ".agents", "plugins", "marketplace.json"), map[string]any{"name": "agents-docs", "plugins": list})
}
func claude(root string, packs map[string][]string) error {
	list := []any{}
	names := make([]string, 0, len(packs))
	for n := range packs {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		list = append(list, map[string]string{"name": n, "source": "./plugins/" + n})
	}
	return write(filepath.Join(root, ".claude-plugin", "marketplace.json"), map[string]any{"name": "agents-docs", "owner": map[string]string{"name": "agents-docs"}, "plugins": list})
}
