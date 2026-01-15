package app

import (
	"os"
	"path/filepath"
	"testing"

	"generator/internal/infra"
)

func TestGeneratorRun(t *testing.T) {
	t.Parallel()

	inRoot := t.TempDir()
	outRoot := t.TempDir()

	writeFile := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	writeFile(filepath.Join(inRoot, "Agents.md"), "agents root")
	writeFile(filepath.Join(inRoot, "agents", "roles", "role.md"), "role doc")
	writeFile(filepath.Join(inRoot, "skills", "skills.md"), "skill index")
	writeFile(filepath.Join(inRoot, "skills", "roles", "scout.md"), "scout skill")

	gen := NewGenerator(infra.OSFS{})
	if err := gen.Run(inRoot, outRoot); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	assertFile := func(path, expected string) {
		t.Helper()
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if string(b) != expected {
			t.Fatalf("content %s = %q, want %q", path, string(b), expected)
		}
	}

	assertFile(filepath.Join(outRoot, "Agents.md"), "agents root")
	assertFile(filepath.Join(outRoot, "agents", "roles", "role.md"), "role doc")
	assertFile(filepath.Join(outRoot, "skills", "roles", "scout", "SKILL.md"), "scout skill\n")

	if _, err := os.Stat(filepath.Join(outRoot, "skills", "skills", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("expected skills/skills/SKILL.md to be excluded, err=%v", err)
	}
}
