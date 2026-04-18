package domain

import (
	"path/filepath"
	"testing"
)

func TestSkillOutputPath(t *testing.T) {
	t.Parallel()

	outRoot := filepath.FromSlash("out/.codex")
	relPath := filepath.FromSlash("skills/roles/scout.md")
	expected := filepath.Join(outRoot, "skills", "roles", "scout", "SKILL.md")
	if got := SkillOutputPath(outRoot, relPath); got != expected {
		t.Fatalf("SkillOutputPath(%q, %q) = %q, want %q", outRoot, relPath, got, expected)
	}
}

func TestFlatSkillOutputPath(t *testing.T) {
	t.Parallel()

	outRoot := filepath.FromSlash("out/.claude")
	expected := filepath.Join("out", ".claude", "skills", "language-go", "SKILL.md")
	if got := FlatSkillOutputPath(outRoot, "language-go"); got != expected {
		t.Fatalf("FlatSkillOutputPath(%q, %q) = %q, want %q", outRoot, "language-go", got, expected)
	}
}
