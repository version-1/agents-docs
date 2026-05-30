package external

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deploy/internal/config"
)

const testCommit = "0123456789abcdef0123456789abcdef01234567"

func TestLoadRequiresCommit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "external-skills.json")
	if err := os.WriteFile(path, []byte(`[
  {"name":"skill","url":"https://github.com/owner/repo/tree/main/skill","type":"git","destination":["dest/skill"]}
]`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected missing commit error")
	}
	if !strings.Contains(err.Error(), "commit is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsInvalidCommit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "external-skills.json")
	if err := os.WriteFile(path, []byte(`[
  {"name":"skill","url":"https://github.com/owner/repo/tree/main/skill","type":"git","commit":"not-a-commit","destination":["dest/skill"]}
]`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected invalid commit error")
	}
	if !strings.Contains(err.Error(), "40-character lowercase hex") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGitFetcherFetchVerifiesCommit(t *testing.T) {
	var calls []string
	fetcher := GitFetcher{runGit: func(args ...string) (string, error) {
		calls = append(calls, strings.Join(args, " "))
		if isRevParseHead(args) {
			return testCommit + "\n", nil
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
}

func TestGitFetcherFetchRejectsCommitMismatch(t *testing.T) {
	actual := "abcdef0123456789abcdef0123456789abcdef01"
	fetcher := GitFetcher{runGit: func(args ...string) (string, error) {
		if isRevParseHead(args) {
			return actual + "\n", nil
		}
		return "", nil
	}}

	_, err := fetcher.Fetch(testSkill("skill"), t.TempDir())
	if err == nil {
		t.Fatal("expected commit mismatch")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("expected %s, got %s", testCommit, actual)) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func isRevParseHead(args []string) bool {
	return len(args) >= 3 && args[len(args)-2] == "rev-parse" && args[len(args)-1] == "HEAD"
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
		Commit:      testCommit,
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
