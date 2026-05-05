package deploy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ExternalSkill struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Destination []string `json:"destination"`
}

type externalSkillFetcher interface {
	Fetch(skill ExternalSkill, workDir string) (string, error)
}

type gitExternalSkillFetcher struct{}

type githubTreeURL struct {
	owner string
	repo  string
	ref   string
	path  string
}

func LoadExternalSkills(path string) ([]ExternalSkill, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read external skills %q: %w", path, err)
	}

	var skills []ExternalSkill
	if err := json.Unmarshal(b, &skills); err != nil {
		return nil, fmt.Errorf("parse external skills %q: %w", path, err)
	}
	for i, skill := range skills {
		if skill.Name == "" {
			return nil, fmt.Errorf("externalSkills[%d].name is required", i)
		}
		if skill.URL == "" {
			return nil, fmt.Errorf("externalSkills[%d].url is required", i)
		}
		if skill.Type == "" {
			return nil, fmt.Errorf("externalSkills[%d].type is required", i)
		}
		if skill.Type != "git" {
			return nil, fmt.Errorf("externalSkills[%d].type %q is not supported", i, skill.Type)
		}
		if len(skill.Destination) == 0 {
			return nil, fmt.Errorf("externalSkills[%d].destination must include at least one path", i)
		}
		for j, destination := range skill.Destination {
			if destination == "" {
				return nil, fmt.Errorf("externalSkills[%d].destination[%d] is required", i, j)
			}
		}
		if _, err := parseGitHubTreeURL(skill.URL); err != nil {
			return nil, fmt.Errorf("externalSkills[%d].url: %w", i, err)
		}
	}
	return skills, nil
}

func validateExternalSkillConflicts(skills []ExternalSkill, cfg Config, cwd string) error {
	names := map[string]string{}
	for i, skill := range skills {
		if previous, ok := names[skill.Name]; ok {
			return fmt.Errorf("external skill name conflict %q from externalSkills[%s] and externalSkills[%d]", skill.Name, previous, i)
		}
		names[skill.Name] = fmt.Sprintf("%d", i)
	}

	internalNames, err := collectInternalSkillNames(cfg, cwd)
	if err != nil {
		return err
	}
	for _, skill := range skills {
		if source, ok := internalNames[skill.Name]; ok {
			return fmt.Errorf("external skill %q conflicts with internal skill at %q", skill.Name, source)
		}
	}

	destinations := map[string]string{}
	for i, skill := range skills {
		for j, destination := range skill.Destination {
			expanded, err := expandHome(destination)
			if err != nil {
				return fmt.Errorf("externalSkills[%d].destination[%d]: %w", i, j, err)
			}
			key := filepath.Clean(expanded)
			if previous, ok := destinations[key]; ok {
				return fmt.Errorf("external skill destination conflict %q from %s and externalSkills[%d].destination[%d]", destination, previous, i, j)
			}
			destinations[key] = fmt.Sprintf("externalSkills[%d].destination[%d]", i, j)
		}
	}
	return nil
}

func collectInternalSkillNames(cfg Config, cwd string) (map[string]string, error) {
	names := map[string]string{}
	for i, item := range cfg.Items {
		src, err := resolveSourcePath(cwd, item.Source)
		if err != nil {
			return nil, fmt.Errorf("resolve source for items[%d]: %w", i, err)
		}
		info, err := os.Stat(src)
		if err != nil {
			return nil, fmt.Errorf("stat source for items[%d] %q: %w", i, src, err)
		}
		if !info.IsDir() {
			continue
		}
		matcher, err := newExcludeMatcher(item.Exclude)
		if err != nil {
			return nil, fmt.Errorf("items[%d]: %w", i, err)
		}
		if err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}
			if matcher.Match(rel) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !d.IsDir() {
				return nil
			}
			skillFile := filepath.Join(path, "SKILL.md")
			skillRel, err := filepath.Rel(src, skillFile)
			if err != nil {
				return err
			}
			if matcher.Match(skillRel) {
				return nil
			}
			info, err := os.Stat(skillFile)
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if info.Mode().IsRegular() {
				names[filepath.Base(path)] = path
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return names, nil
}

func validateFetchedExternalSkill(skill ExternalSkill, src string) error {
	skillFile := filepath.Join(src, "SKILL.md")
	info, err := os.Stat(skillFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("external skill %q does not contain SKILL.md at %q", skill.Name, src)
		}
		return fmt.Errorf("stat external skill %q SKILL.md: %w", skill.Name, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("external skill %q SKILL.md is not a regular file", skill.Name)
	}
	actualName, err := readSkillName(skillFile)
	if err != nil {
		return fmt.Errorf("read external skill %q name: %w", skill.Name, err)
	}
	if actualName != "" && actualName != skill.Name {
		return fmt.Errorf("external skill name mismatch: config name %q but SKILL.md name is %q", skill.Name, actualName)
	}
	return nil
}

func readSkillName(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", nil
	}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "---" {
			return "", nil
		}
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "name:")), nil
		}
	}
	return "", nil
}

func (gitExternalSkillFetcher) Fetch(skill ExternalSkill, workDir string) (string, error) {
	treeURL, err := parseGitHubTreeURL(skill.URL)
	if err != nil {
		return "", err
	}

	repoDir := filepath.Join(workDir, safePathName(skill.Name))
	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", treeURL.owner, treeURL.repo)
	if err := runGit("clone", "--depth", "1", "--filter=blob:none", "--sparse", "--branch", treeURL.ref, cloneURL, repoDir); err != nil {
		return "", err
	}
	if err := runGit("-C", repoDir, "sparse-checkout", "set", "--", treeURL.path); err != nil {
		return "", err
	}
	return filepath.Join(repoDir, filepath.FromSlash(treeURL.path)), nil
}

func parseGitHubTreeURL(raw string) (githubTreeURL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return githubTreeURL{}, fmt.Errorf("parse GitHub tree URL: %w", err)
	}
	if u.Scheme != "https" || u.Host != "github.com" {
		return githubTreeURL{}, fmt.Errorf("only https://github.com/<owner>/<repo>/tree/<ref>/<path> URLs are supported")
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 5 || parts[2] != "tree" {
		return githubTreeURL{}, fmt.Errorf("expected https://github.com/<owner>/<repo>/tree/<ref>/<path>")
	}
	owner, repo, ref := parts[0], parts[1], parts[3]
	path := strings.Join(parts[4:], "/")
	if owner == "" || repo == "" || ref == "" || path == "" {
		return githubTreeURL{}, fmt.Errorf("expected non-empty owner, repo, ref, and path")
	}
	return githubTreeURL{owner: owner, repo: repo, ref: ref, path: path}, nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return nil
}

func safePathName(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_")
	return replacer.Replace(name)
}
