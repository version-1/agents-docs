package external

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deploy/internal/config"
)

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
