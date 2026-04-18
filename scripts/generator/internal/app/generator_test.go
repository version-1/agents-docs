package app

import (
	"os"
	"path/filepath"
	"testing"

	"generator/internal/infra"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertFile(t *testing.T, path, expected string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(b) != expected {
		t.Fatalf("content %s = %q, want %q", path, string(b), expected)
	}
}

func assertNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to not exist, err=%v", path, err)
	}
}

func TestGeneratorRun_Nested(t *testing.T) {
	t.Parallel()

	inRoot := t.TempDir()
	outRoot := t.TempDir()

	writeFile(t, filepath.Join(inRoot, "Agents.md"), "agents root")
	writeFile(t, filepath.Join(inRoot, "agents", "roles", "role.md"), "role doc")
	writeFile(t, filepath.Join(inRoot, "skills", "skills.md"), "skill index")
	writeFile(t, filepath.Join(inRoot, "skills", "roles", "scout.md"), "---\nname: scout\n---\nscout skill")

	gen := NewGenerator(infra.OSFS{}, NewPathRespectSkillGenerator())
	if err := gen.Run(inRoot, outRoot); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	assertFile(t, filepath.Join(outRoot, "Agents.md"), "agents root")
	assertFile(t, filepath.Join(outRoot, "agents", "roles", "role.md"), "role doc")
	assertFile(t, filepath.Join(outRoot, "skills", "roles", "scout", "SKILL.md"), "---\nname: scout\n---\nscout skill\n")

	assertNotExist(t, filepath.Join(outRoot, "skills", "skills", "SKILL.md"))
}

func TestGeneratorRun_FlatSkills(t *testing.T) {
	t.Parallel()

	inRoot := t.TempDir()
	outRoot := t.TempDir()

	writeFile(t, filepath.Join(inRoot, "Agents.md"), "agents root")
	writeFile(t, filepath.Join(inRoot, "agents", "roles", "role.md"), "role doc")
	writeFile(t, filepath.Join(inRoot, "skills", "skills.md"), "skill index")
	writeFile(t, filepath.Join(inRoot, "skills", "roles", "scout.md"), "---\nname: role-scout\n---\nscout skill")
	writeFile(t, filepath.Join(inRoot, "skills", "languages", "go", "go.md"), "---\nname: language-go\n---\ngo skill")

	gen := NewGenerator(infra.OSFS{}, NewFlatSkillGenerator())
	if err := gen.Run(inRoot, outRoot); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	// flat mode: フロントマターの name がディレクトリ名になる
	assertFile(t, filepath.Join(outRoot, "skills", "role-scout", "SKILL.md"), "---\nname: role-scout\n---\nscout skill\n")
	assertFile(t, filepath.Join(outRoot, "skills", "language-go", "SKILL.md"), "---\nname: language-go\n---\ngo skill\n")

	// agents はそのままコピー
	assertFile(t, filepath.Join(outRoot, "Agents.md"), "agents root")
	assertFile(t, filepath.Join(outRoot, "agents", "roles", "role.md"), "role doc")

	// nested パスには存在しない
	assertNotExist(t, filepath.Join(outRoot, "skills", "roles", "scout", "SKILL.md"))
	assertNotExist(t, filepath.Join(outRoot, "skills", "languages", "go", "go", "SKILL.md"))
}

func TestGeneratorRun_FlatSkills_MissingName(t *testing.T) {
	t.Parallel()

	inRoot := t.TempDir()
	outRoot := t.TempDir()

	writeFile(t, filepath.Join(inRoot, "Agents.md"), "agents root")
	writeFile(t, filepath.Join(inRoot, "agents", "roles", "role.md"), "role doc")
	writeFile(t, filepath.Join(inRoot, "skills", "bad.md"), "no frontmatter here")

	gen := NewGenerator(infra.OSFS{}, NewFlatSkillGenerator())
	if err := gen.Run(inRoot, outRoot); err == nil {
		t.Fatal("expected error for missing frontmatter name, got nil")
	}
}
