package external

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deploy/internal/config"
)

const testTreeHash = "0123456789abcdef0123456789abcdef01234567"

func TestLoadRequiresTreeHash(t *testing.T) {
	path := filepath.Join(t.TempDir(), "external-skills.json")
	if err := os.WriteFile(path, []byte(`[
  {"name":"skill","url":"https://github.com/owner/repo/tree/main/skill","type":"git","destination":["dest/skill"]}
]`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected missing treeHash error")
	}
	if !strings.Contains(err.Error(), "treeHash is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsInvalidTreeHash(t *testing.T) {
	path := filepath.Join(t.TempDir(), "external-skills.json")
	if err := os.WriteFile(path, []byte(`[
  {"name":"skill","url":"https://github.com/owner/repo/tree/main/skill","type":"git","treeHash":"not-a-tree-hash","destination":["dest/skill"]}
]`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected invalid treeHash error")
	}
	if !strings.Contains(err.Error(), "40-character lowercase hex") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAcceptsGitHubSkillBlobURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "external-skills.json")
	if err := os.WriteFile(path, []byte(`[
  {"name":"grilling","url":"https://github.com/owner/repo/blob/main/skills/grilling/SKILL.md","type":"git","treeHash":"0123456789abcdef0123456789abcdef01234567","destination":["dest/grilling"]}
]`), 0644); err != nil {
		t.Fatal(err)
	}

	skills, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "grilling" {
		t.Fatalf("unexpected skills: %#v", skills)
	}
}

func TestLoadRejectsGitHubBlobURLForNonSkillFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "external-skills.json")
	if err := os.WriteFile(path, []byte(`[
  {"name":"skill","url":"https://github.com/owner/repo/blob/main/skill/README.md","type":"git","treeHash":"0123456789abcdef0123456789abcdef01234567","destination":["dest/skill"]}
]`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected invalid blob URL error")
	}
	if !strings.Contains(err.Error(), "must end in SKILL.md") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGitFetcherFetchVerifiesTreeHash(t *testing.T) {
	var calls []string
	fetcher := GitFetcher{runGit: func(args ...string) (string, error) {
		calls = append(calls, strings.Join(args, " "))
		if isRevParseTree(args, "skill") {
			return testTreeHash + "\n", nil
		}
		return "", nil
	}}

	src, err := fetcher.Fetch(testSkill("skill"), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(src, filepath.Join("skill")) {
		t.Fatalf("unexpected source path: %s", src)
	}
	if len(calls) != 3 {
		t.Fatalf("expected clone, rev-parse, sparse-checkout calls, got %v", calls)
	}
	if isRevParseHead(calls) {
		t.Fatalf("fetch should verify the skill tree hash, not repository HEAD: %v", calls)
	}
}

func TestGitFetcherFetchUsesParentDirectoryForSkillBlobURL(t *testing.T) {
	var calls []string
	skill := testSkill("grilling")
	skill.URL = "https://github.com/owner/repo/blob/main/skills/productivity/grilling/SKILL.md"
	fetcher := GitFetcher{runGit: func(args ...string) (string, error) {
		calls = append(calls, strings.Join(args, " "))
		if isRevParseTree(args, "skills/productivity/grilling") {
			return testTreeHash + "\n", nil
		}
		return "", nil
	}}

	src, err := fetcher.Fetch(skill, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(src, filepath.Join("skills", "productivity", "grilling")) {
		t.Fatalf("unexpected source path: %s", src)
	}
	if !strings.HasSuffix(calls[2], "sparse-checkout set -- skills/productivity/grilling") {
		t.Fatalf("unexpected sparse-checkout call: %v", calls)
	}
}

func TestGitFetcherFetchChecksOutCommitRef(t *testing.T) {
	var calls []string
	skill := testSkill("grilling")
	skill.URL = "https://github.com/owner/repo/blob/697d4ce9742da558fd1ba6697c8e9775e2e302dd/skills/productivity/grilling/SKILL.md"
	fetcher := GitFetcher{runGit: func(args ...string) (string, error) {
		calls = append(calls, strings.Join(args, " "))
		if isRevParseTree(args, "skills/productivity/grilling") {
			return testTreeHash + "\n", nil
		}
		return "", nil
	}}

	if _, err := fetcher.Fetch(skill, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 5 {
		t.Fatalf("expected clone, fetch, checkout, rev-parse, sparse-checkout calls, got %v", calls)
	}
	if !strings.Contains(calls[0], "--no-checkout") || strings.Contains(calls[0], "--branch") {
		t.Fatalf("commit ref clone should not use --branch: %v", calls)
	}
	if !strings.HasSuffix(calls[1], "fetch --depth 1 origin 697d4ce9742da558fd1ba6697c8e9775e2e302dd") {
		t.Fatalf("unexpected commit fetch call: %v", calls)
	}
	if !strings.HasSuffix(calls[2], "checkout --detach FETCH_HEAD") {
		t.Fatalf("unexpected commit checkout call: %v", calls)
	}
}

func TestGitFetcherFetchRejectsTreeHashMismatch(t *testing.T) {
	actual := "abcdef0123456789abcdef0123456789abcdef01"
	fetcher := GitFetcher{runGit: func(args ...string) (string, error) {
		if isRevParseTree(args, "skill") {
			return actual + "\n", nil
		}
		return "", nil
	}}

	_, err := fetcher.Fetch(testSkill("skill"), t.TempDir())
	if err == nil {
		t.Fatal("expected tree hash mismatch")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("expected %s, got %s", testTreeHash, actual)) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func isRevParseTree(args []string, path string) bool {
	return len(args) >= 3 && args[len(args)-2] == "rev-parse" && args[len(args)-1] == "HEAD:"+path
}

func isRevParseHead(calls []string) bool {
	for _, call := range calls {
		if strings.HasSuffix(call, " rev-parse HEAD") {
			return true
		}
	}
	return false
}

func TestValidateConflictsRejectsInternalSkillNameConflict(t *testing.T) {
	root := t.TempDir()
	writeTestSkill(t, filepath.Join(root, "codex", "skills", "coding"))

	err := ValidateConflicts([]Skill{testSkill("coding")}, config.Config{
		Items: []config.Item{{Source: "codex/skills", Destination: "dest"}},
	}, root)
	if err == nil {
		t.Fatal("expected internal skill conflict")
	}
	if !strings.Contains(err.Error(), "conflicts with internal skill") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConflictsIgnoresExcludedInternalSkillDir(t *testing.T) {
	root := t.TempDir()
	writeTestSkill(t, filepath.Join(root, "codex", "skills", "coding"))

	err := ValidateConflicts([]Skill{testSkill("coding")}, config.Config{
		Items: []config.Item{{
			Source:      "codex/skills",
			Destination: "dest",
			Exclude:     []string{"coding"},
		}},
	}, root)
	if err != nil {
		t.Fatalf("excluded internal skill should not conflict: %v", err)
	}
}

func TestValidateConflictsIgnoresInternalSkillWithExcludedSkillFile(t *testing.T) {
	root := t.TempDir()
	writeTestSkill(t, filepath.Join(root, "codex", "skills", "coding"))

	err := ValidateConflicts([]Skill{testSkill("coding")}, config.Config{
		Items: []config.Item{{
			Source:      "codex/skills",
			Destination: "dest",
			Exclude:     []string{"coding/SKILL.md"},
		}},
	}, root)
	if err != nil {
		t.Fatalf("internal skill with excluded SKILL.md should not conflict: %v", err)
	}
}

func TestInternalSkillScanKeyNormalizesExcludePatterns(t *testing.T) {
	a := internalSkillScanKey("/repo/skills", []string{"./a/", "b"})
	b := internalSkillScanKey("/repo/skills", []string{"b", "a"})
	if a != b {
		t.Fatalf("equivalent excludes should produce same key: %q != %q", a, b)
	}

	c := internalSkillScanKey("/repo/other", []string{"b", "a"})
	if a == c {
		t.Fatalf("different source should produce different key: %q", a)
	}
}

func testSkill(name string) Skill {
	return Skill{
		Name:        name,
		URL:         "https://github.com/owner/repo/tree/main/" + name,
		Type:        "git",
		TreeHash:    testTreeHash,
		Destination: []string{"dest/" + name},
	}
}

func writeTestSkill(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: test\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
}
