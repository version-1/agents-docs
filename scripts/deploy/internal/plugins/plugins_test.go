package plugins

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deploy/internal/external"
)

type recordingFetcher struct {
	sources map[string]string
	calls   []string
}

func (f *recordingFetcher) Fetch(s external.Skill, _ string) (string, error) {
	f.calls = append(f.calls, s.Name)
	p, ok := f.sources[s.Name]
	if !ok {
		return "", errors.New("unexpected fetch: " + s.Name)
	}
	return p, nil
}

func TestBuildGeneratesAllGroupAndIndividualPlugins(t *testing.T) {
	r := fixture(t, map[string][]string{"code": {"code-general", "code-go"}, "docs": {"documenting"}})
	for _, n := range []string{"code-general", "code-go", "documenting"} {
		writeSkill(t, filepath.Join(r, "skills", "internal", n), n)
	}
	out := filepath.Join(r, "nested", "marketplace")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	for _, plugin := range []string{"code", "docs", "skill-code-general", "skill-code-go", "skill-documenting"} {
		assertExists(t, filepath.Join(out, "plugins", plugin, ".codex-plugin", "plugin.json"))
		assertExists(t, filepath.Join(out, "plugins", plugin, ".claude-plugin", "plugin.json"))
	}
	assertExists(t, filepath.Join(out, ".agents", "plugins", "marketplace.json"))
	assertExists(t, filepath.Join(out, ".claude-plugin", "marketplace.json"))
}

func TestBuildResolvesIncludedGroupsAndDeduplicatesDiamondMembership(t *testing.T) {
	r := fixtureGroups(t, map[string]Group{
		"alpha":  {Skills: []string{"a"}},
		"beta":   {Skills: []string{"b"}},
		"middle": {Includes: []string{"alpha", "beta", "alpha"}},
		"top":    {Includes: []string{"middle", "beta"}},
	})
	for _, name := range []string{"a", "b"} {
		writeSkill(t, filepath.Join(r, "skills", "internal", name), name)
	}
	out := filepath.Join(r, "out")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	for _, plugin := range []string{"alpha", "beta", "middle", "top"} {
		assertExists(t, filepath.Join(out, "plugins", plugin, ".codex-plugin", "plugin.json"))
	}
	for _, name := range []string{"a", "b"} {
		assertExists(t, filepath.Join(out, "plugins", "top", "skills", name, "SKILL.md"))
	}
	entries, err := os.ReadDir(filepath.Join(out, "plugins", "top", "skills"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("top skill count = %d, want 2", len(entries))
	}
}

func TestBuildRejectsInvalidIncludesWithoutReplacingOutput(t *testing.T) {
	for _, tc := range []struct {
		name   string
		groups map[string]Group
		want   string
	}{
		{"unknown", map[string]Group{"one": {Includes: []string{"missing"}}}, "unknown included group"},
		{"invalid name", map[string]Group{"one": {Includes: []string{"Bad"}}}, "invalid included group"},
		{"self cycle", map[string]Group{"one": {Includes: []string{"one"}}}, "cycle"},
		{"mutual cycle", map[string]Group{"one": {Includes: []string{"two"}}, "two": {Includes: []string{"one"}}}, "cycle"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := fixtureGroups(t, tc.groups)
			out := filepath.Join(r, "out")
			mustWrite(t, filepath.Join(out, "keep"), "keep")
			_, err := Build(options(r, out))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
			assertFile(t, filepath.Join(out, "keep"), "keep")
		})
	}
}

func TestBuildIncludedGroupMembershipChangesCompositeVersion(t *testing.T) {
	r := fixtureGroups(t, map[string]Group{
		"base": {Skills: []string{"alpha"}},
		"core": {Includes: []string{"base"}},
	})
	skill := filepath.Join(r, "skills", "internal", "alpha")
	writeSkill(t, skill, "alpha")
	first := filepath.Join(r, "first")
	if _, err := Build(options(r, first)); err != nil {
		t.Fatal(err)
	}
	writeSkill(t, filepath.Join(r, "skills", "internal", "beta"), "beta")
	writeJSON(t, filepath.Join(r, "skill-groups.json"), GroupsFile{Groups: map[string]Group{
		"base": {Skills: []string{"alpha", "beta"}},
		"core": {Includes: []string{"base"}},
	}})
	second := filepath.Join(r, "second")
	if _, err := Build(options(r, second)); err != nil {
		t.Fatal(err)
	}
	if version(t, first, "core") == version(t, second, "core") {
		t.Fatal("composite version did not change with included group membership")
	}
}

func TestBuildRejectsInvalidCatalogWithoutReplacingOutput(t *testing.T) {
	for _, tc := range []struct {
		name   string
		groups map[string][]string
		want   string
	}{
		{"invalid group", map[string][]string{"Bad": {"alpha"}}, "invalid group"}, {"reserved group", map[string][]string{"skill-alpha": {"alpha"}}, "reserved"}, {"invalid skill", map[string][]string{"one": {"../alpha"}}, "invalid skill"}, {"duplicate skill", map[string][]string{"one": {"alpha"}, "two": {"alpha"}}, "registered"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := fixture(t, tc.groups)
			writeSkill(t, filepath.Join(r, "skills", "internal", "alpha"), "alpha")
			out := filepath.Join(r, "out")
			mustWrite(t, filepath.Join(out, "keep"), "keep")
			_, err := Build(options(r, out))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
			assertFile(t, filepath.Join(out, "keep"), "keep")
		})
	}
}

func TestBuildFetchesAllExternalAndRejectsNameConflict(t *testing.T) {
	r := fixture(t, map[string][]string{"one": {"alpha"}})
	writeSkill(t, filepath.Join(r, "skills", "internal", "alpha"), "alpha")
	remote := filepath.Join(r, "remote")
	writeSkill(t, remote, "remote")
	writeExternal(t, r, "remote")
	f := &recordingFetcher{sources: map[string]string{"remote": remote}}
	if _, err := Build(Options{GroupsPath: filepath.Join(r, "skill-groups.json"), ExternalPath: filepath.Join(r, "external-skills.json"), Output: filepath.Join(r, "out"), Fetcher: f}); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(f.calls, ","); got != "remote" {
		t.Fatalf("fetches = %q", got)
	}
	assertExists(t, filepath.Join(r, "out", "plugins", "skill-remote", "skills", "remote", "SKILL.md"))
	writeSkill(t, filepath.Join(r, "skills", "internal", "remote"), "remote")
	if _, err := Build(Options{GroupsPath: filepath.Join(r, "skill-groups.json"), ExternalPath: filepath.Join(r, "external-skills.json"), Output: filepath.Join(r, "conflict")}); err == nil || !strings.Contains(err.Error(), "conflict") {
		t.Fatalf("conflict error = %v", err)
	}
}

func TestBuildIndividualPluginIncludesRelativeSkillClosure(t *testing.T) {
	r := fixture(t, map[string][]string{"code": {"code-general", "code-go"}})
	writeSkill(t, filepath.Join(r, "skills", "internal", "code-general"), "code-general")
	writeSkill(t, filepath.Join(r, "skills", "internal", "code-go"), "code-go")
	mustWrite(t, filepath.Join(r, "skills", "internal", "code-go", "SKILL.md"), "---\nname: code-go\n---\nRead ../code-general/SKILL.md\n")
	out := filepath.Join(r, "out")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(out, "plugins", "skill-code-go", "skills", "code-general", "SKILL.md"))
}

func TestBuildIndividualPluginIncludesRelativeSkillResourceClosure(t *testing.T) {
	r := fixture(t, map[string][]string{"code": {"alpha", "beta"}})
	writeSkill(t, filepath.Join(r, "skills", "internal", "alpha"), "alpha")
	writeSkill(t, filepath.Join(r, "skills", "internal", "beta"), "beta")
	mustWrite(t, filepath.Join(r, "skills", "internal", "alpha", "SKILL.md"), "---\nname: alpha\n---\nRead ../beta/references/lang.md\n")
	mustWrite(t, filepath.Join(r, "skills", "internal", "beta", "references", "lang.md"), "resource\n")
	out := filepath.Join(r, "out")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(out, "plugins", "skill-alpha", "skills", "beta", "references", "lang.md"))
}

func TestBuildOutputSafetyModeAndVersion(t *testing.T) {
	r := fixture(t, map[string][]string{"one": {"alpha"}})
	skill := filepath.Join(r, "skills", "internal", "alpha")
	writeSkill(t, skill, "alpha")
	mustWrite(t, filepath.Join(skill, "run.sh"), "#!/bin/sh\n")
	if err := os.Chmod(filepath.Join(skill, "run.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(r, "out")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(filepath.Join(out, "plugins", "skill-alpha", "skills", "alpha", "run.sh")); err != nil || info.Mode()&0111 == 0 {
		t.Fatalf("mode=%v err=%v", info, err)
	}
	v := version(t, out, "skill-alpha")
	if _, err := Build(options(r, filepath.Join(r, "two"))); err != nil {
		t.Fatal(err)
	}
	if v != version(t, filepath.Join(r, "two"), "skill-alpha") {
		t.Fatal("unstable version")
	}
	mustWrite(t, filepath.Join(skill, "SKILL.md"), "---\nname: alpha\n---\nchanged\n")
	if _, err := Build(options(r, filepath.Join(r, "three"))); err != nil {
		t.Fatal(err)
	}
	if v == version(t, filepath.Join(r, "three"), "skill-alpha") {
		t.Fatal("version did not change")
	}
	unsafe := filepath.Join(r, "unsafe")
	mustWrite(t, filepath.Join(unsafe, "keep"), "keep")
	if _, err := Build(options(r, unsafe)); err == nil {
		t.Fatal("expected unmarked output error")
	}
	assertFile(t, filepath.Join(unsafe, "keep"), "keep")
}

func TestBuildPreservesUsersPreviousSibling(t *testing.T) {
	r := fixture(t, map[string][]string{"one": {"alpha"}})
	writeSkill(t, filepath.Join(r, "skills", "internal", "alpha"), "alpha")
	out := filepath.Join(r, "out")
	mustWrite(t, filepath.Join(out+".previous", "user-directory", "keep"), "keep")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	assertFile(t, filepath.Join(out+".previous", "user-directory", "keep"), "keep")
}

func TestBuildIndividualClosuresAreTransitive(t *testing.T) {
	r := fixture(t, map[string][]string{"code": {"a", "b", "c"}})
	for _, name := range []string{"a", "b", "c"} {
		writeSkill(t, filepath.Join(r, "skills", "internal", name), name)
	}
	mustWrite(t, filepath.Join(r, "skills", "internal", "a", "SKILL.md"), "---\nname: a\n---\nRead ../b/SKILL.md\n")
	mustWrite(t, filepath.Join(r, "skills", "internal", "b", "SKILL.md"), "---\nname: b\n---\nRead ../c/SKILL.md\n")
	out := filepath.Join(r, "out")
	if _, err := Build(options(r, out)); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a", "b", "c"} {
		assertExists(t, filepath.Join(out, "plugins", "skill-a", "skills", name, "SKILL.md"))
	}
}

func TestBuildIndividualPluginVersionIncludesClosureAndFileMetadata(t *testing.T) {
	r := fixture(t, map[string][]string{"code": {"a", "b"}})
	writeSkill(t, filepath.Join(r, "skills", "internal", "a"), "a")
	writeSkill(t, filepath.Join(r, "skills", "internal", "b"), "b")
	mustWrite(t, filepath.Join(r, "skills", "internal", "a", "SKILL.md"), "---\nname: a\n---\nRead ../b/SKILL.md\n")
	build := func(name string) string {
		out := filepath.Join(r, name)
		if _, err := Build(options(r, out)); err != nil {
			t.Fatal(err)
		}
		return version(t, out, "skill-a")
	}
	first := build("one")
	mustWrite(t, filepath.Join(r, "skills", "internal", "b", "SKILL.md"), "---\nname: b\n---\nchanged\n")
	if got := build("two"); got == first {
		t.Fatal("closure content did not affect version")
	}
	mustWrite(t, filepath.Join(r, "skills", "internal", "a", "original.txt"), "same")
	beforeRename := build("three")
	if err := os.Rename(filepath.Join(r, "skills", "internal", "a", "original.txt"), filepath.Join(r, "skills", "internal", "a", "renamed.txt")); err != nil {
		t.Fatal(err)
	}
	second := build("four")
	if second == beforeRename {
		t.Fatal("filename change did not affect version")
	}
	if got := build("same-tree"); got != second {
		t.Fatal("same tree produced unstable version")
	}
	if err := os.Chmod(filepath.Join(r, "skills", "internal", "a", "renamed.txt"), 0755); err != nil {
		t.Fatal(err)
	}
	if got := build("five"); got == second {
		t.Fatal("mode change did not affect version")
	}
}

func TestBuildMarketplaceIsByteDeterministic(t *testing.T) {
	r := fixture(t, map[string][]string{"zeta": {"z"}, "alpha": {"a"}})
	writeSkill(t, filepath.Join(r, "skills", "internal", "a"), "a")
	writeSkill(t, filepath.Join(r, "skills", "internal", "z"), "z")
	first, second := filepath.Join(r, "one"), filepath.Join(r, "two")
	if _, err := Build(options(r, first)); err != nil {
		t.Fatal(err)
	}
	if _, err := Build(options(r, second)); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{filepath.Join(".agents", "plugins", "marketplace.json"), filepath.Join(".claude-plugin", "marketplace.json"), "summary.json"} {
		one, err := os.ReadFile(filepath.Join(first, file))
		if err != nil {
			t.Fatal(err)
		}
		two, err := os.ReadFile(filepath.Join(second, file))
		if err != nil {
			t.Fatal(err)
		}
		if string(one) != string(two) {
			t.Fatalf("%s differs between identical builds", file)
		}
	}
}

func TestRepositoryCatalogBuildsOffline(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", "..", "..", ".."))
	var groups GroupsFile
	if err := load(filepath.Join(repo, "skill-groups.json"), &groups); err != nil {
		t.Fatal(err)
	}
	var ext []external.Skill
	if err := load(filepath.Join(repo, "external-skills.json"), &ext); err != nil {
		t.Fatal(err)
	}
	expectedSkills := repositorySourceSkills(t, repo, ext)
	groupedSkills := map[string]struct{}{}
	for _, group := range groups.Groups {
		for _, skill := range group.Skills {
			groupedSkills[skill] = struct{}{}
		}
	}
	assertSameNames(t, "configured skills", groupedSkills, expectedSkills)
	sources := map[string]string{}
	for _, s := range ext {
		p := filepath.Join(t.TempDir(), s.Name)
		writeSkill(t, p, s.Name)
		sources[s.Name] = p
	}
	out := filepath.Join(t.TempDir(), "out")
	if _, err := Build(Options{GroupsPath: filepath.Join(repo, "skill-groups.json"), ExternalPath: filepath.Join(repo, "external-skills.json"), Output: out, Fetcher: &recordingFetcher{sources: sources}}); err != nil {
		t.Fatal(err)
	}
	expectedPlugins := map[string]struct{}{}
	for name := range groups.Groups {
		expectedPlugins[name] = struct{}{}
	}
	for skill := range expectedSkills {
		expectedPlugins["skill-"+skill] = struct{}{}
		assertPluginManifests(t, out, "skill-"+skill)
		assertExists(t, filepath.Join(out, "plugins", "skill-"+skill, "skills", skill, "SKILL.md"))
	}
	for group := range groups.Groups {
		assertPluginManifests(t, out, group)
	}
	assertSameNames(t, "generated plugins", directoryNames(t, filepath.Join(out, "plugins")), expectedPlugins)
}

func repositorySourceSkills(t *testing.T, repo string, externalSkills []external.Skill) map[string]struct{} {
	t.Helper()
	skills := map[string]struct{}{}
	entries, err := os.ReadDir(filepath.Join(repo, "skills", "internal"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		info, err := os.Stat(filepath.Join(repo, "skills", "internal", entry.Name(), "SKILL.md"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if !info.Mode().IsRegular() {
			continue
		}
		addName(t, skills, entry.Name(), "internal source")
	}
	for _, skill := range externalSkills {
		addName(t, skills, skill.Name, "external config")
	}
	return skills
}

func addName(t *testing.T, names map[string]struct{}, name, source string) {
	t.Helper()
	if _, exists := names[name]; exists {
		t.Fatalf("duplicate skill %q in %s", name, source)
	}
	names[name] = struct{}{}
}

func directoryNames(t *testing.T, path string) map[string]struct{} {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	names := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names[entry.Name()] = struct{}{}
		}
	}
	return names
}

func assertSameNames(t *testing.T, description string, got, want map[string]struct{}) {
	t.Helper()
	for name := range want {
		if _, exists := got[name]; !exists {
			t.Errorf("%s missing %q", description, name)
		}
	}
	for name := range got {
		if _, exists := want[name]; !exists {
			t.Errorf("%s unexpectedly contains %q", description, name)
		}
	}
}

func assertPluginManifests(t *testing.T, output, plugin string) {
	t.Helper()
	assertExists(t, filepath.Join(output, "plugins", plugin, ".codex-plugin", "plugin.json"))
	assertExists(t, filepath.Join(output, "plugins", plugin, ".claude-plugin", "plugin.json"))
}

func fixture(t *testing.T, groups map[string][]string) string {
	t.Helper()
	g := GroupsFile{Groups: map[string]Group{}}
	for n, s := range groups {
		g.Groups[n] = Group{Skills: s}
	}
	return fixtureGroupsFile(t, g)
}
func fixtureGroups(t *testing.T, groups map[string]Group) string {
	t.Helper()
	return fixtureGroupsFile(t, GroupsFile{Groups: groups})
}
func fixtureGroupsFile(t *testing.T, groups GroupsFile) string {
	t.Helper()
	r := t.TempDir()
	writeJSON(t, filepath.Join(r, "skill-groups.json"), groups)
	writeJSON(t, filepath.Join(r, "external-skills.json"), []external.Skill{})
	return r
}
func options(root, out string) Options {
	return Options{GroupsPath: filepath.Join(root, "skill-groups.json"), ExternalPath: filepath.Join(root, "external-skills.json"), Output: out}
}
func writeSkill(t *testing.T, d, n string) {
	t.Helper()
	mustWrite(t, filepath.Join(d, "SKILL.md"), "---\nname: "+n+"\n---\nbody\n")
}
func writeExternal(t *testing.T, r, n string) {
	t.Helper()
	writeJSON(t, filepath.Join(r, "external-skills.json"), []external.Skill{{Name: n, URL: "https://github.com/owner/repo/tree/main/" + n, Type: "git", TreeHash: "0123456789abcdef0123456789abcdef01234567"}})
}
func writeJSON(t *testing.T, p string, v any) {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	mustWrite(t, p, string(b))
}
func mustWrite(t *testing.T, p, s string) {
	t.Helper()
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(p, []byte(s), 0644); e != nil {
		t.Fatal(e)
	}
}
func assertExists(t *testing.T, p string) {
	t.Helper()
	if _, e := os.Stat(p); e != nil {
		t.Fatalf("missing %s: %v", p, e)
	}
}
func assertFile(t *testing.T, p, want string) {
	t.Helper()
	b, e := os.ReadFile(p)
	if e != nil || string(b) != want {
		t.Fatalf("%s=%q (%v)", p, b, e)
	}
}
func version(t *testing.T, root, plugin string) string {
	t.Helper()
	var m struct {
		Version string `json:"version"`
	}
	b, e := os.ReadFile(filepath.Join(root, "plugins", plugin, ".codex-plugin", "plugin.json"))
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(b, &m); e != nil {
		t.Fatal(e)
	}
	return m.Version
}
